#!/usr/bin/env perl
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This program reads a file containing function prototypes
# (like syscall_darwin.go) and generates system call bodies.
# The prototypes are marked by lines beginning with "//sys"
# and read like func declarations if //sys is replaced by func, but:
#	* The parameter lists must give a name for each argument.
#	  This includes return parameters.
#	* The parameter lists must give a type for each argument:
#	  the (x, y, z int) shorthand is not allowed.
#	* If the return parameter is an error number, it must be named errno.

# A line beginning with //sysnb is like //sys, except that the
# goroutine will not be suspended during the execution of the system
# call.  This must only be used for system calls which can never
# block, as otherwise the system call could cause all goroutines to
# hang.

use strict;

my $cmdline = "mksyscall.pl " . join(' ', @ARGV);
my $errors = 0;
my $_32bit = "";
my $plan9 = 0;
my $openbsd = 0;
my $netbsd = 0;
my $dragonfly = 0;
my $arm = 0; # 64-bit value should use (even, odd)-pair
my $tags = "";  # build tags

if($ARGV[0] eq "-b32") {
	$_32bit = "big-endian";
	shift;
} elsif($ARGV[0] eq "-l32") {
	$_32bit = "little-endian";
	shift;
}
if($ARGV[0] eq "-plan9") {
	$plan9 = 1;
	shift;
}
if($ARGV[0] eq "-openbsd") {
	$openbsd = 1;
	shift;
}
if($ARGV[0] eq "-netbsd") {
	$netbsd = 1;
	shift;
}
if($ARGV[0] eq "-dragonfly") {
	$dragonfly = 1;
	shift;
}
if($ARGV[0] eq "-arm") {
	$arm = 1;
	shift;
}
if($ARGV[0] eq "-tags") {
	shift;
	$tags = $ARGV[0];
	shift;
}

if($ARGV[0] =~ /^-/) {
	print STDERR "usage: mksyscall.pl [-b32 | -l32] [-tags x,y] [file ...]\n";
	exit 1;
}

# Check that we are using the new build system if we should
if($ENV{'GOOS'} eq "linux" && $ENV{'GOARCH'} ne "sparc64") {
	if($ENV{'GOLANG_SYS_BUILD'} ne "docker") {
		print STDERR "In the new build system, mksyscall should not be called directly.\n";
		print STDERR "See README.md\n";
		exit 1;
	}
}


sub parseparamlist($) {
	my ($list) = @_;
	$list =~ s/^\s*//;
	$list =~ s/\s*$//;
	if($list eq "") {
		return ();
	}
	return split(/\s*,\s*/, $list);
}

sub parseparam($) {
	my ($p) = @_;
	if($p !~ /^(\S*) (\S*)$/) {
		print STDERR "$ARGV:$.: malformed parameter: $p\n";
		$errors = 1;
		return ("xx", "int");
	}
	return ($1, $2);
}

my $text = "";
while(<>) {
	chomp;
	s/\s+/ /g;
	s/^\s+//;
	s/\s+$//;
	my $nonblock = /^\/\/sysnb /;
	next if !/^\/\/sys / && !$nonblock;

	# Line must be of the form
	#	func Open(path string, mode int, perm int) (fd int, errno error)
	# Split into name, in params, out params.
	if(!/^\/\/sys(nb)? (\w+)\(([^()]*)\)\s*(?:\(([^()]+)\))?\s*(?:=\s*((?i)SYS_[A-Z0-9_]+))?$/) {
		print STDERR "$ARGV:$.: malformed //sys declaration\n";
		$errors = 1;
		next;
	}
	my ($func, $in, $out, $sysname) = ($2, $3, $4, $5);

	# Split argument lists on comma.
	my @in = parseparamlist($in);
	my @out = parseparamlist($out);

	# Try in vain to keep people from editing this file.
	# The theory is that they jump into the middle of the file
	# without reading the header.
	$text .= "// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT\n\n";

	# Go function header.
	my $out_decl = @out ? sprintf(" (%s)", join(', ', @out)) : "";
	$text .= sprintf "func %s(%s)%s {\n", $func, join(', ', @in), $out_decl;

	# Check if err return available
	my $errvar = "";
	foreach my $p (@out) {
		my ($name, $type) = parseparam($p);
		if($type eq "error") {
			$errvar = $name;
			last;
		}
	}

	# Prepare arguments to Syscall.
	my @args = ();
	my $n = 0;
	foreach my $p (@in) {
		my ($name, $type) = parseparam($p);
		if($type =~ /^\*/) {
			push @args, "uintptr(unsafe.Pointer($name))";
		} elsif($type eq "string" && $errvar ne "") {
			$text .= "\tvar _p$n *byte\n";
			$text .= "\t_p$n, $errvar = BytePtrFromString($name)\n";
			$text .= "\tif $errvar != nil {\n\t\treturn\n\t}\n";
			push @args, "uintptr(unsafe.Pointer(_p$n))";
			$n++;
		} elsif($type eq "string") {
			print STDERR "$ARGV:$.: $func uses string arguments, but has no error return\n";
			$text .= "\tvar _p$n *byte\n";
			$text .= "\t_p$n, _ = BytePtrFromString($name)\n";
			push @args, "uintptr(unsafe.Pointer(_p$n))";
			$n++;
		} elsif($type =~ /^\[\](.*)/) {
			# Convert slice into pointer, length.
			# Have to be careful not to take address of &a[0] if len == 0:
			# pass dummy pointer in that case.
			# Used to pass nil, but some OSes or simulators reject write(fd, nil, 0).
			$text .= "\tvar _p$n unsafe.Pointer\n";
			$text .= "\tif len($name) > 0 {\n\t\t_p$n = unsafe.Pointer(\&${name}[0])\n\t}";
			$text .= " else {\n\t\t_p$n = unsafe.Pointer(&_zero)\n\t}";
			$text .= "\n";
			push @args, "uintptr(_p$n)", "uintptr(len($name))";
			$n++;
		} elsif($type eq "int64" && ($openbsd || $netbsd)) {
			push @args, "0";
			if($_32bit eq "big-endian") {
				push @args, "uintptr($name>>32)", "uintptr($name)";
			} elsif($_32bit eq "little-endian") {
				push @args, "uintptr($name)", "uintptr($name>>32)";
			} else {
				push @args, "uintptr($name)";
			}
		} elsif($type eq "int64" && $dragonfly) {
			if ($func !~ /^extp(read|write)/i) {
				push @args, "0";
			}
			if($_32bit eq "big-endian") {
				push @args, "uintptr($name>>32)", "uintptr($name)";
			} elsif($_32bit eq "little-endian") {
				push @args, "uintptr($name)", "uintptr($name>>32)";
			} else {
				push @args, "uintptr($name)";
			}
		} elsif($type eq "int64" && $_32bit ne "") {
			if(@args % 2 && $arm) {
				# arm abi specifies 64-bit argument uses
				# (even, odd) pair
				push @args, "0"
			}
			if($_32bit eq "big-endian") {
				push @args, "uintptr($name>>32)", "uintptr($name)";
			} else {
				push @args, "uintptr($name)", "uintptr($name>>32)";
			}
		} else {
			push @args, "uintptr($name)";
		}
	}

	# Determine which form to use; pad args with zeros.
	my $asm = "Syscall";
	if ($nonblock) {
		if ($errvar ne "") {
			$asm = "RawSyscall";
		} else {
			$asm = "RawSyscallNoError";
		}
	} else {
		if ($errvar eq "") {
			$asm = "SyscallNoError";
		}
	}
	if(@args <= 3) {
		while(@args < 3) {
			push @args, "0";
		}
	} elsif(@args <= 6) {
		$asm .= "6";
		while(@args < 6) {
			push @args, "0";
		}
	} elsif(@args <= 9) {
		$asm .= "9";
		while(@args < 9) {
			push @args, "0";
		}
	} else {
		print STDERR "$ARGV:$.: too many arguments to system call\n";
	}

	# System call number.
	if($sysname eq "") {
		$sysname = "SYS_$func";
		$sysname =~ s/([a-z])([A-Z])/${1}_$2/g;	# turn FooBar into Foo_Bar
		$sysname =~ y/a-z/A-Z/;
	}

	# Actual call.
	my $args = join(', ', @args);
	my $call = "$asm($sysname, $args)";

	# Assign return values.
	my $body = "";
	my @ret = ("_", "_", "_");
	my $do_errno = 0;
	for(my $i=0; $i<@out; $i++) {
		my $p = $out[$i];
		my ($name, $type) = parseparam($p);
		my $reg = "";
		if($name eq "err" && !$plan9) {
			$reg = "e1";
			$ret[2] = $reg;
			$do_errno = 1;
		} elsif($name eq "err" && $plan9) {
			$ret[0] = "r0";
			$ret[2] = "e1";
			next;
		} else {
			$reg = sprintf("r%d", $i);
			$ret[$i] = $reg;
		}
		if($type eq "bool") {
			$reg = "$reg != 0";
		}
		if($type eq "int64" && $_32bit ne "") {
			# 64-bit number in r1:r0 or r0:r1.
			if($i+2 > @out) {
				print STDERR "$ARGV:$.: not enough registers for int64 return\n";
			}
			if($_32bit eq "big-endian") {
				$reg = sprintf("int64(r%d)<<32 | int64(r%d)", $i, $i+1);
			} else {
				$reg = sprintf("int64(r%d)<<32 | int64(r%d)", $i+1, $i);
			}
			$ret[$i] = sprintf("r%d", $i);
			$ret[$i+1] = sprintf("r%d", $i+1);
		}
		if($reg ne "e1" || $plan9) {
			$body .= "\t$name = $type($reg)\n";
		}
	}
	if ($ret[0] eq "_" && $ret[1] eq "_" && $ret[2] eq "_") {
		$text .= "\t$call\n";
	} else {
		if ($errvar ne "") {
			$text .= "\t$ret[0], $ret[1], $ret[2] := $call\n";
		} else {
			$text .= "\t$ret[0], $ret[1] := $call\n";
		}
	}
	$text .= $body;

	if ($plan9 && $ret[2] eq "e1") {
		$text .= "\tif int32(r0) == -1 {\n";
		$text .= "\t\terr = e1\n";
		$text .= "\t}\n";
	} elsif ($do_errno) {
		$text .= "\tif e1 != 0 {\n";
		$text .= "\t\terr = errnoErr(e1)\n";
		$text .= "\t}\n";
	}
	$text .= "\treturn\n";
	$text .= "}\n\n";
}

chomp $text;
chomp $text;

if($errors) {
	exit 1;
}

print <<EOF;
// $cmdline
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build $tags

package unix

import (
	"syscall"
	"unsafe"
)

var _ syscall.Errno

$text
EOF
exit 0;

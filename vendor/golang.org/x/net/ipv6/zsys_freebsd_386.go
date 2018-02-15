// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_freebsd.go

package ipv6

const (
	sysIPV6_UNICAST_HOPS   = 0x4
	sysIPV6_MULTICAST_IF   = 0x9
	sysIPV6_MULTICAST_HOPS = 0xa
	sysIPV6_MULTICAST_LOOP = 0xb
	sysIPV6_JOIN_GROUP     = 0xc
	sysIPV6_LEAVE_GROUP    = 0xd
	sysIPV6_PORTRANGE      = 0xe
	sysICMP6_FILTER        = 0x12

	sysIPV6_CHECKSUM = 0x1a
	sysIPV6_V6ONLY   = 0x1b

	sysIPV6_IPSEC_POLICY = 0x1c

	sysIPV6_RTHDRDSTOPTS = 0x23

	sysIPV6_RECVPKTINFO  = 0x24
	sysIPV6_RECVHOPLIMIT = 0x25
	sysIPV6_RECVRTHDR    = 0x26
	sysIPV6_RECVHOPOPTS  = 0x27
	sysIPV6_RECVDSTOPTS  = 0x28

	sysIPV6_USE_MIN_MTU = 0x2a
	sysIPV6_RECVPATHMTU = 0x2b

	sysIPV6_PATHMTU = 0x2c

	sysIPV6_PKTINFO  = 0x2e
	sysIPV6_HOPLIMIT = 0x2f
	sysIPV6_NEXTHOP  = 0x30
	sysIPV6_HOPOPTS  = 0x31
	sysIPV6_DSTOPTS  = 0x32
	sysIPV6_RTHDR    = 0x33

	sysIPV6_RECVTCLASS = 0x39

	sysIPV6_AUTOFLOWLABEL = 0x3b

	sysIPV6_TCLASS   = 0x3d
	sysIPV6_DONTFRAG = 0x3e

	sysIPV6_PREFER_TEMPADDR = 0x3f

	sysIPV6_BINDANY = 0x40

	sysIPV6_MSFILTER = 0x4a

	sysMCAST_JOIN_GROUP         = 0x50
	sysMCAST_LEAVE_GROUP        = 0x51
	sysMCAST_JOIN_SOURCE_GROUP  = 0x52
	sysMCAST_LEAVE_SOURCE_GROUP = 0x53
	sysMCAST_BLOCK_SOURCE       = 0x54
	sysMCAST_UNBLOCK_SOURCE     = 0x55

	sysIPV6_PORTRANGE_DEFAULT = 0x0
	sysIPV6_PORTRANGE_HIGH    = 0x1
	sysIPV6_PORTRANGE_LOW     = 0x2

	sizeofSockaddrStorage = 0x80
	sizeofSockaddrInet6   = 0x1c
	sizeofInet6Pktinfo    = 0x14
	sizeofIPv6Mtuinfo     = 0x20

	sizeofIPv6Mreq       = 0x14
	sizeofGroupReq       = 0x84
	sizeofGroupSourceReq = 0x104

	sizeofICMPv6Filter = 0x20
)

type sockaddrStorage struct {
	Len         uint8
	Family      uint8
	X__ss_pad1  [6]int8
	X__ss_align int64
	X__ss_pad2  [112]int8
}

type sockaddrInet6 struct {
	Len      uint8
	Family   uint8
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type inet6Pktinfo struct {
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
}

type ipv6Mtuinfo struct {
	Addr sockaddrInet6
	Mtu  uint32
}

type ipv6Mreq struct {
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
}

type groupReq struct {
	Interface uint32
	Group     sockaddrStorage
}

type groupSourceReq struct {
	Interface uint32
	Group     sockaddrStorage
	Source    sockaddrStorage
}

type icmpv6Filter struct {
	Filt [8]uint32
}

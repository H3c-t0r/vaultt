// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

var logger hclog.Logger

func fatal(err error) {
	logger.Error("fatal error", "error", err)
	os.Exit(1)
}

type generator struct {
	file *ast.File
	fset *token.FileSet
}

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:  "stubmaker",
		Level: hclog.Trace,
	})

	// Setup git, both so we can determine if we're running on enterprise, and
	// so we can make sure we don't clobber a non-transient file.
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		if err.Error() != "repository does not exist" {
			fatal(err)
		}
		repo = nil
	}

	var wt *git.Worktree
	if repo != nil {
		wt, err = repo.Worktree()
		if err != nil {
			fatal(err)
		}
		if !isEnterprise(wt) {
			return
		}
	}

	// Read the file and figure out if we need to do anything.
	inputFile := os.Getenv("GOFILE")
	if !strings.HasSuffix(inputFile, "_stubs_oss.go") {
		fatal(fmt.Errorf("stubmaker should only be invoked from files ending in _stubs_oss.go"))
	}

	baseFilename := strings.TrimSuffix(inputFile, "_stubs_oss.go")
	outputFile := baseFilename + "_stubs_ent.go"
	b, err := os.ReadFile(inputFile)
	if err != nil {
		fatal(err)
	}

	inputParsed, err := parseFile(b)
	if err != nil {
		fatal(err)
	}
	needed, existing, err := inputParsed.areStubsNeeded()
	if err != nil {
		fatal(err)
	}
	if !needed {
		return
	}

	// We'd like to write the file, but first make sure that we're not going
	// to blow away anyone's work or overwrite a file already in git.
	if repo != nil {
		head, err := repo.Head()
		if err != nil {
			fatal(err)
		}
		obj, err := repo.Object(plumbing.AnyObject, head.Hash())
		if err != nil {
			fatal(err)
		}

		st, err := wt.Status()
		if err != nil {
			fatal(err)
		}

		tracked, err := inGit(wt, st, obj, outputFile)
		if err != nil {
			fatal(err)
		}
		if tracked {
			fatal(fmt.Errorf("output file %s exists in git, not overwriting", outputFile))
		}
	}

	// Now we can finally write the file
	output, err := os.Create(outputFile + ".tmp")
	if err != nil {
		fatal(err)
	}
	err = inputParsed.writeStubs(output, existing)
	if err != nil {
		// If we don't end up writing to the file, delete it.
		os.Remove(outputFile + ".tmp")
	} else {
		os.Rename(outputFile+".tmp", outputFile)
	}
	if err != nil {
		fatal(err)
	}
}

func (g *generator) writeStubs(output *os.File, existingFuncs map[string]struct{}) error {
	// delete all functions/methods that are already defined
	g.modifyAST(existingFuncs)

	// write the updated code to buf
	buf := new(bytes.Buffer)
	err := format.Node(buf, g.fset, g.file)
	if err != nil {
		return err
	}

	// remove any unneeded imports
	res, err := imports.Process("", buf.Bytes(), &imports.Options{
		Fragment:   true,
		AllErrors:  false,
		Comments:   true,
		FormatOnly: false,
	})
	if err != nil {
		return err
	}

	// add the code generation line and update the build tags
	outputLines, err := fixGeneratedComments(res)
	if err != nil {
		return err
	}
	_, err = output.WriteString(strings.Join(outputLines, "\n") + "\n")
	return err
}

func fixGeneratedComments(b []byte) ([]string, error) {
	warning := "// Code generated by tools/stubmaker; DO NOT EDIT."
	goGenerate := "//go:generate go run github.com/hashicorp/vault/tools/stubmaker"

	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	var outputLines []string
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.Contains(line, "//go:build ") && strings.Contains(line, "!enterprise"):
			outputLines = append(outputLines, warning, "")
			line = strings.ReplaceAll(line, "!enterprise", "enterprise")
		case line == goGenerate:
			continue
		}
		outputLines = append(outputLines, line)
	}
	return outputLines, scanner.Err()
}

func inGit(wt *git.Worktree, st git.Status, obj object.Object, path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("path %s can't be made absolute: %w", path, err)
	}
	relPath, err := filepath.Rel(wt.Filesystem.Root(), absPath)
	if err != nil {
		return false, fmt.Errorf("path %s can't be made relative: %w", absPath, err)
	}

	fst := st.File(relPath)
	if fst.Worktree != git.Untracked || fst.Staging != git.Untracked {
		return true, nil
	}

	curwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	blob, err := resolve(obj, relPath)
	if err != nil && !strings.Contains(err.Error(), "file not found") {
		return false, fmt.Errorf("error resolving path %s from %s: %w", relPath, curwd, err)
	}

	return blob != nil, nil
}

func isEnterprise(wt *git.Worktree) bool {
	st, err := wt.Filesystem.Stat("enthelpers")
	onOss := errors.Is(err, os.ErrNotExist)
	onEnt := st != nil

	switch {
	case onOss && !onEnt:
	case !onOss && onEnt:
	default:
		fatal(err)
	}
	return onEnt
}

// resolve blob at given path from obj. obj can be a commit, tag, tree, or blob.
func resolve(obj object.Object, path string) (*object.Blob, error) {
	switch o := obj.(type) {
	case *object.Commit:
		t, err := o.Tree()
		if err != nil {
			return nil, err
		}
		return resolve(t, path)
	case *object.Tag:
		target, err := o.Object()
		if err != nil {
			return nil, err
		}
		return resolve(target, path)
	case *object.Tree:
		file, err := o.File(path)
		if err != nil {
			return nil, err
		}
		return &file.Blob, nil
	case *object.Blob:
		return o, nil
	default:
		return nil, object.ErrUnsupportedObject
	}
}

// areStubsNeeded checks if all functions and methods defined in the stub file
// are present in the package
func (g *generator) areStubsNeeded() (needed bool, existingStubs map[string]struct{}, err error) {
	pkg, err := parsePackage(".", []string{"enterprise"})
	if err != nil {
		return false, nil, err
	}

	stubFunctions := make(map[string]struct{})
	for _, d := range g.file.Decls {
		dFunc, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		stubFunctions[dFunc.Name.Name] = struct{}{}

	}
	found := make(map[string]struct{})
	for name, val := range pkg.TypesInfo.Defs {
		if val == nil {
			continue
		}
		_, ok := val.Type().(*types.Signature)
		if !ok {
			continue
		}
		if _, ok := stubFunctions[name.Name]; ok {
			found[name.Name] = struct{}{}
		}
	}

	return len(found) != len(stubFunctions), found, nil
}

func (g *generator) modifyAST(exists map[string]struct{}) {
	astutil.Apply(g.file, nil, func(c *astutil.Cursor) bool {
		switch x := c.Node().(type) {
		case *ast.FuncDecl:
			if _, ok := exists[x.Name.Name]; ok {
				c.Delete()
			}
		}

		return true
	})
}

func parsePackage(name string, tags []string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, name)
	if err != nil {
		return nil, fmt.Errorf("error parsing package %s: %v", name, err)
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("error: %d packages found", len(pkgs))
	}
	return pkgs[0], nil
}

func parseFile(buffer []byte) (*generator, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", buffer, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return &generator{
		file: f,
		fset: fs,
	}, nil
}

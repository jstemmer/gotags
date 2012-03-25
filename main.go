package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

const VERSION = "0.0.1"

var fset *token.FileSet

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(1)
	}

	filename := os.Args[1]

	fset = token.NewFileSet()
	tags := make(sort.StringSlice, 0)

	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// header
	tags = append(tags, "!_TAG_FILE_FORMAT\t2\t")
	tags = append(tags, "!_TAG_FILE_SORTED\t1\t")

	// package
	if f.Name != nil {
		tags = append(tags, createTag(f.Name.Name, f.Name.Pos(), "p"))
	}

	// imports
	for _, im := range f.Imports {
		if im.Path != nil {
			name := strings.Trim(im.Path.Value, "\"")
			tags = append(tags, createTag(name, im.Path.Pos(), "i"))
		}
	}

	// declarations
	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			if decl.Name != nil {
				// TODO: add params, receiver, etc
				tags = append(tags, createTag(decl.Name.Name, decl.Pos(), "f"))
			}
		case *ast.GenDecl:
			for _, s := range decl.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					if ts.Name != nil {
						tags = append(tags, createTag(ts.Name.Name, ts.Pos(), "s"))
					}
				}
			}
		}
	}

	// sort and print tags
	sort.Sort(tags)
	for _, tag := range tags {
		fmt.Println(tag)
	}
}

func printUsage() {
	fmt.Printf("gotags version %s\n\n", VERSION)
	fmt.Printf("Usage: %s file\n", os.Args[0])
}

func createTag(name string, pos token.Pos, tagtype string) string {
	return NewTag(name, fset.File(pos).Name(), fset.Position(pos).Line, tagtype).String()
}

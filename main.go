package main

import (
	"bytes"
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
		// TODO: fix error handling; it should still result in a valid ctags file
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// header
	tags = append(tags, "!_TAG_FILE_FORMAT\t2\t")
	tags = append(tags, "!_TAG_FILE_SORTED\t1\t")

	// package
	if f.Name != nil {
		tags = append(tags, createTag(f.Name.Name, f.Name.Pos(), "p").String())
	}

	// imports
	for _, im := range f.Imports {
		if im.Path != nil {
			name := strings.Trim(im.Path.Value, "\"")
			tags = append(tags, createTag(name, im.Path.Pos(), "i").String())
		}
	}

	// declarations
	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			tags = append(tags, createFuncTag(decl))
		case *ast.GenDecl:
			for _, s := range decl.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					tags = append(tags, createTag(ts.Name.Name, ts.Pos(), "t").String())
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

func createTag(name string, pos token.Pos, tagtype string) *Tag {
	return NewTag(name, fset.File(pos).Name(), fset.Position(pos).Line, tagtype)
}

func createFuncTag(f *ast.FuncDecl) string {
	if f == nil || f.Name == nil {
		return ""
	}

	tag := createTag(f.Name.Name, f.Pos(), "f")

	// access
	if ast.IsExported(tag.Name) {
		tag.Fields["access"] = "public"
	} else {
		tag.Fields["access"] = "private"
	}

	// signature
	var sig bytes.Buffer
	sig.WriteByte('(')
	for i, param := range f.Type.Params.List {
		// parameter names
		for j, n := range param.Names {
			sig.WriteString(n.Name)
			if j < len(param.Names)-1 {
				sig.WriteString(", ")
			}
		}

		// parameter type
		sig.WriteByte(' ')
		sig.WriteString(getType(param.Type))

		if i < len(f.Type.Params.List)-1 {
			sig.WriteString(", ")
		}
	}
	sig.WriteByte(')')
	tag.Fields["signature"] = sig.String()

	// receiver
	if f.Recv != nil && len(f.Recv.List) > 0 {
		tag.Fields["type"] = getType(f.Recv.List[0].Type)
	}

	return tag.String()
}

func getType(node ast.Node) (paramType string) {
	switch t := node.(type) {
	case *ast.Ident:
		paramType = t.Name
	case *ast.StarExpr:
		paramType = "*" + getType(t.X)
	case *ast.SelectorExpr:
		paramType = getType(t.X) + "." + getType(t.Sel)
	}
	return
}

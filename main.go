package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
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

	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// header
	fmt.Println("!_TAG_FILE_FORMAT\t2\t")
	fmt.Println("!_TAG_FILE_SORTED\t0\t")

	// package
	if f.Name != nil {
		printTag(f.Name.Name, "p", f.Name.Pos())
	}

	// imports
	for _, im := range f.Imports {
		if im.Path != nil {
			printTag(strings.Trim(im.Path.Value, "\""), "i", im.Path.Pos())
		}
	}
}

func printUsage() {
	fmt.Printf("gotags version %s\n\n", VERSION)
	fmt.Printf("Usage: %s file\n", os.Args[0])
}

func printTag(tag, kind string, pos token.Pos) {
	line := fset.Position(pos).Line
	file := fset.File(pos).Name()
	fmt.Printf("%s\t%s\t%d;\"\t%s\n", tag, file, line, kind)
}

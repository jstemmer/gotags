package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const VERSION = "0.0.1"

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(1)
	}

	filename := os.Args[1]

	fset := token.NewFileSet()
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
		line := fset.Position(f.Name.Pos()).Line
		name := fset.File(f.Name.Pos()).Name()
		fmt.Printf("%s\t%s\t%d;\"\tp\n", f.Name.Name, name, line)
	}

	// imports
	for _, im := range f.Imports {
		if im.Path != nil {
			line := fset.Position(im.Path.Pos()).Line
			name := fset.File(im.Path.Pos()).Name()
			fmt.Printf("%s\t%s\t%d;\"\ti\n", strings.Trim(im.Path.Value, "\""), name, line)
		}
	}
}

func printUsage() {
	fmt.Printf("gotags version %s\n\n", VERSION)
	fmt.Printf("Usage: %s file\n", os.Args[0])
}

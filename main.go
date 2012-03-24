package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	filename := os.Args[1]

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// header
	fmt.Println("!_TAG_FILE_FORMAT\t2\t//")
	fmt.Println("!_TAG_FILE_SORTED\t0\t//")

	// package
	if f.Name != nil {
		line := fset.Position(f.Name.Pos()).Line
		fmt.Printf("%s\t%s\t%d;\" p\n", f.Name.Name, filename, line)
	}
}

func usage() {
	// TODO: print usage
	os.Exit(1)
}

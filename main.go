package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
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
	fmt.Println("!_TAG_FILE_FORMAT\t2\t//")
	fmt.Println("!_TAG_FILE_SORTED\t0\t//")

	// package
	if f.Name != nil {
		line := fset.Position(f.Name.Pos()).Line
		fmt.Printf("%s\t%s\t%d;\" p\n", f.Name.Name, filename, line)
	}
}

func printUsage() {
	fmt.Printf("gotags version %s\n\n", VERSION)
	fmt.Printf("Usage: %s file\n", os.Args[0])
}

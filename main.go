package main

import (
	"fmt"
	"os"
)

const VERSION = "0.0.1"

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(1)
	}

	tags, err := Parse(os.Args[1])
	if err != nil {
		// TODO: fix error handling; it should still result in a valid ctags file
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// TODO: sort

	// header
	fmt.Println("!_TAG_FILE_FORMAT\t2\t")
	fmt.Println("!_TAG_FILE_SORTED\t0\t")
	for _, tag := range tags {
		fmt.Println(tag.String())
	}
}

func printUsage() {
	fmt.Printf("gotags version %s\n\n", VERSION)
	fmt.Printf("Usage: %s file\n", os.Args[0])
}

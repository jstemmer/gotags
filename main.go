package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

const VERSION = "0.0.1"

var (
	sortOutput bool
	silent     bool
	printTree  bool
)

func init() {
	flag.BoolVar(&sortOutput, "sort", true, "sort tags")
	flag.BoolVar(&silent, "silent", false, "do not produce any output on error")
	flag.BoolVar(&printTree, "tree", false, "print syntax tree (debugging)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gotags version %s\n\n", VERSION)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file(s)\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "no file specified\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if printTree {
		PrintTree(flag.Arg(0))
		return
	}

	tags := make([]Tag, 0)
	for _, file := range flag.Args() {
		ts, err := Parse(file)
		if err != nil {
			if !silent {
				fmt.Fprintf(os.Stderr, "parse error: %s\n\n", err)
			}
			continue
		}
		tags = append(tags, ts...)
	}

	output := createMetaTags()
	for _, tag := range tags {
		output = append(output, tag.String())
	}

	if sortOutput {
		sort.Sort(sort.StringSlice(output))
	}

	for _, s := range output {
		fmt.Println(s)
	}
}

func createMetaTags() []string {
	var sorted int
	if sortOutput {
		sorted = 1
	}
	return []string{
		"!_TAG_FILE_FORMAT\t2\t",
		fmt.Sprintf("!_TAG_FILE_SORTED\t%d\t", sorted),
		"!_TAG_PROGRAM_NAME\tgotags\t",
		"!_TAG_PROGRAM_URL\thttps://github.com/jstemmer/gotags\t",
		fmt.Sprintf("!_TAG_PROGRAM_VERSION\t%s\t", VERSION),
	}
}

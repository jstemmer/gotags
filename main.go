package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

const (
	VERSION      = "1.2.0"
	NAME         = "gotags"
	URL          = "https://github.com/jstemmer/gotags"
	AUTHOR_NAME  = "Joel Stemmer"
	AUTHOR_EMAIL = "stemmertech@gmail.com"
)

var (
	printVersion bool
	sortOutput   bool
	silent       bool
	printTree    bool // for debugging
)

// Initialize flags.
func init() {
	flag.BoolVar(&printVersion, "v", false, "print version")
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

	if printVersion {
		fmt.Printf("gotags version %s\n", VERSION)
		return
	}

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "no file specified\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if printTree {
		PrintTree(flag.Arg(0))
		return
	}

	tags := []Tag{}
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

// createMetaTags returns a list of meta tags.
func createMetaTags() []string {
	var sorted int
	if sortOutput {
		sorted = 1
	}
	return []string{
		"!_TAG_FILE_FORMAT\t2\t",
		fmt.Sprintf("!_TAG_FILE_SORTED\t%d\t/0=unsorted, 1=sorted/", sorted),
		fmt.Sprintf("!_TAG_PROGRAM_AUTHOR\t%s\t/%s/", AUTHOR_NAME, AUTHOR_EMAIL),
		fmt.Sprintf("!_TAG_PROGRAM_NAME\t%s\t", NAME),
		fmt.Sprintf("!_TAG_PROGRAM_URL\t%s\t", URL),
		fmt.Sprintf("!_TAG_PROGRAM_VERSION\t%s\t", VERSION),
	}
}

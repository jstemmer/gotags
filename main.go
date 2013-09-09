package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Contants used for the meta tags
const (
	Version     = "1.3.0"
	Name        = "gotags"
	URL         = "https://github.com/jstemmer/gotags"
	AuthorName  = "Joel Stemmer"
	AuthorEmail = "stemmertech@gmail.com"
)

var (
	printVersion bool
	inputFile    string
	recurse      bool
	sortOutput   bool
	silent       bool
)

// Initialize flags.
func init() {
	flag.BoolVar(&printVersion, "v", false, "print version")
	flag.StringVar(&inputFile, "L", "", "source file names are read from the specified file.")
	flag.BoolVar(&recurse, "R", false, "recurse into directories in the file list")
	flag.BoolVar(&sortOutput, "sort", true, "sort tags")
	flag.BoolVar(&silent, "silent", false, "do not produce any output on error")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gotags version %s\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file(s)\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func walkDir(names []string, dir string) ([]string, error) {
	e := filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			names = append(names, path)
		}
		return nil
	})

	return names, e
}

func recurseNames(names []string) ([]string, error) {
	var ret []string
	for _, name := range names {
		info, e := os.Stat(name)
		if e != nil || info == nil || !info.IsDir() {
			ret = append(ret, name) // defer the error handling to the scanner
		} else {
			ret, e = walkDir(ret, name)
			if e != nil {
				return names, e
			}
		}
	}
	return ret, nil
}

func readNames(names []string) ([]string, error) {
	if len(inputFile) == 0 {
		return names, nil
	}

	var scanner *bufio.Scanner
	if inputFile != "-" {
		in, err := os.Open(inputFile)
		if err != nil {
			return nil, err
		}

		defer in.Close()
		scanner = bufio.NewScanner(in)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func getFileNames() ([]string, error) {
	var names []string

	names = append(names, flag.Args()...)
	names, err := readNames(names)
	if err != nil {
		return nil, err
	}

	if recurse {
		names, err = recurseNames(names)
		if err != nil {
			return nil, err
		}
	}

	return names, nil
}

func main() {
	flag.Parse()

	if printVersion {
		fmt.Printf("gotags version %s\n", Version)
		return
	}

	files, err := getFileNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get specified files\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if len(files) == 0 && len(inputFile) == 0 {
		fmt.Fprintf(os.Stderr, "no file specified\n\n")
		flag.Usage()
		os.Exit(1)
	}

	tags := []Tag{}
	for _, file := range files {
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
		"!_TAG_FILE_FORMAT\t2",
		fmt.Sprintf("!_TAG_FILE_SORTED\t%d\t/0=unsorted, 1=sorted/", sorted),
		fmt.Sprintf("!_TAG_PROGRAM_AUTHOR\t%s\t/%s/", AuthorName, AuthorEmail),
		fmt.Sprintf("!_TAG_PROGRAM_NAME\t%s", Name),
		fmt.Sprintf("!_TAG_PROGRAM_URL\t%s", URL),
		fmt.Sprintf("!_TAG_PROGRAM_VERSION\t%s", Version),
	}
}

package main

import (
	"reflect"
	"strings"
	"testing"
)

var sources = []string{
	"LICENSE",
	"README.md",
	"fields.go",
	"fields_test.go",
	"files_test.go",
	"main.go",
	"parser.go",
	"parser_test.go",
	"tag.go",
	"tag_test.go",
	"tags",
	"tests/const.go-src",
	"tests/func.go-src",
	"tests/import.go-src",
	"tests/interface.go-src",
	"tests/range.go-src",
	"tests/simple.go-src",
	"tests/struct.go-src",
	"tests/type.go-src",
	"tests/var.go-src",
}

func TestCommandLineFiles(t *testing.T) {
	patterns := patternList{}

	files, err := getFileNames(sources, false, patterns)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(sources, files) {
		t.Errorf("%+v != %+v", sources, files)
	}
}

func TestSingleExcludePattern(t *testing.T) {
	// Single pattern - exclude *_test.go files
	patterns := []string{"*_test.go"}
	files, err := getFileNames(sources, false, patterns)
	if err != nil {
		t.Error(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			t.Errorf("%v should not be included", f)
		}
	}
}

func TestRecursiveExcludes(t *testing.T) {
	input := []string{"/usr/local/go/src"}
	patterns := []string{
		"*_test.go",
		"/usr/local/go/src/*/*/testdata/*",
		"/usr/local/go/src/*/*/testdata/*/*",
		"/usr/local/go/src/*/*/testdata/*/*/*",
		"/usr/local/go/src/*/*/testdata/*/*/*/*",
		"/usr/local/go/src/*/*/testdata/*/*/*/*/*",
		"/usr/local/go/src/*/*/testdata/*/*/*/*/*/*",
		"/usr/local/go/src/*/*/testdata/*/*/*/*/*/*/*",
		"/usr/local/go/src/*/*/*/testdata/*",
	}
	files, err := getFileNames(input, true, patterns)
	if err != nil {
		t.Error(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") || strings.Contains(f, "/testdata/") {
			t.Errorf("%v should not be included", f)
		}
	}
	t.Log(files)
}

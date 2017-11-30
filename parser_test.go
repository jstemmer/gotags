package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"testing"
)

var goVersionRegexp = regexp.MustCompile(`^go1(?:\.(\d+))?`)

// This type is used to implement the sort.Interface interface
// in order to be able to sort an array of Tag
type TagSlice []Tag

// Return the len of the array
func (t TagSlice) Len() int {
	return len(t)
}

// Compare two elements of a tag array
func (t TagSlice) Less(i, j int) bool {
	return t[i].String() < t[j].String()
}

// Swap two elements of the underlying array
func (t TagSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Dump the names of the tags in a TagSlice
func (t TagSlice) Dump() {
	for idx, val := range t {
		fmt.Println(idx, val.Name)
	}
}

type F map[TagField]string

var testCases = []struct {
	filename         string
	relative         bool
	basepath         string
	minversion       int
	withExtraSymbols bool
	tags             []Tag
}{
	{filename: "testdata/const.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Constant", 3, "c", F{"access": "public", "type": "string"}),
		tag("OtherConst", 4, "c", F{"access": "public"}),
		tag("A", 7, "c", F{"access": "public"}),
		tag("B", 8, "c", F{"access": "public"}),
		tag("C", 8, "c", F{"access": "public"}),
		tag("D", 9, "c", F{"access": "public"}),
	}},
	{filename: "testdata/const.go", withExtraSymbols: true, tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Constant", 3, "c", F{"access": "public", "type": "string"}),
		tag("OtherConst", 4, "c", F{"access": "public"}),
		tag("A", 7, "c", F{"access": "public"}),
		tag("B", 8, "c", F{"access": "public"}),
		tag("C", 8, "c", F{"access": "public"}),
		tag("D", 9, "c", F{"access": "public"}),
		tag("Test.Constant", 3, "c", F{"access": "public", "type": "string"}),
		tag("Test.OtherConst", 4, "c", F{"access": "public"}),
		tag("Test.A", 7, "c", F{"access": "public"}),
		tag("Test.B", 8, "c", F{"access": "public"}),
		tag("Test.C", 8, "c", F{"access": "public"}),
		tag("Test.D", 9, "c", F{"access": "public"}),
	}},
	{filename: "testdata/func.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Function1", 3, "f", F{"access": "public", "signature": "()", "type": "string"}),
		tag("function2", 6, "f", F{"access": "private", "signature": "(p1, p2 int, p3 *string)"}),
		tag("function3", 9, "f", F{"access": "private", "signature": "()", "type": "bool"}),
		tag("function4", 12, "f", F{"access": "private", "signature": "(p interface{})", "type": "interface{}"}),
		tag("function5", 15, "f", F{"access": "private", "signature": "()", "type": "string, string, error"}),
		tag("function6", 18, "f", F{"access": "private", "signature": "(v ...interface{})"}),
		tag("function7", 21, "f", F{"access": "private", "signature": "(s ...string)"}),
	}},
	{filename: "testdata/func.go", withExtraSymbols: true, tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Test.Function1", 3, "f", F{"access": "public", "signature": "()", "type": "string"}),
		tag("Test.function2", 6, "f", F{"access": "private", "signature": "(p1, p2 int, p3 *string)"}),
		tag("Test.function3", 9, "f", F{"access": "private", "signature": "()", "type": "bool"}),
		tag("Test.function4", 12, "f", F{"access": "private", "signature": "(p interface{})", "type": "interface{}"}),
		tag("Test.function5", 15, "f", F{"access": "private", "signature": "()", "type": "string, string, error"}),
		tag("Test.function6", 18, "f", F{"access": "private", "signature": "(v ...interface{})"}),
		tag("Test.function7", 21, "f", F{"access": "private", "signature": "(s ...string)"}),
		tag("Function1", 3, "f", F{"access": "public", "signature": "()", "type": "string"}),
		tag("function2", 6, "f", F{"access": "private", "signature": "(p1, p2 int, p3 *string)"}),
		tag("function3", 9, "f", F{"access": "private", "signature": "()", "type": "bool"}),
		tag("function4", 12, "f", F{"access": "private", "signature": "(p interface{})", "type": "interface{}"}),
		tag("function5", 15, "f", F{"access": "private", "signature": "()", "type": "string, string, error"}),
		tag("function6", 18, "f", F{"access": "private", "signature": "(v ...interface{})"}),
		tag("function7", 21, "f", F{"access": "private", "signature": "(s ...string)"}),
	}},
	{filename: "testdata/import.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("fmt", 3, "i", F{}),
		tag("go/ast", 6, "i", F{}),
		tag("go/parser", 7, "i", F{}),
	}},
	{filename: "testdata/import.go", withExtraSymbols: true, tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("fmt", 3, "i", F{}),
		tag("go/ast", 6, "i", F{}),
		tag("go/parser", 7, "i", F{}),
	}},
	{filename: "testdata/interface.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("InterfaceMethod", 4, "m", F{"access": "public", "signature": "(int)", "ntype": "Interface", "type": "string"}),
		tag("OtherMethod", 5, "m", F{"access": "public", "signature": "()", "ntype": "Interface"}),
		tag("io.Reader", 6, "e", F{"access": "public", "ntype": "Interface"}),
		tag("Interface", 3, "n", F{"access": "public", "type": "interface"}),
	}},
	{filename: "testdata/interface.go", withExtraSymbols: true, tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("InterfaceMethod", 4, "m", F{"access": "public", "signature": "(int)", "ntype": "Interface", "type": "string"}),
		tag("OtherMethod", 5, "m", F{"access": "public", "signature": "()", "ntype": "Interface"}),
		tag("io.Reader", 6, "e", F{"access": "public", "ntype": "Interface"}),
		tag("Interface", 3, "n", F{"access": "public", "type": "interface"}),
		tag("Test.Interface", 3, "n", F{"access": "public", "type": "interface"}),
	}},
	{filename: "testdata/struct.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Field1", 4, "w", F{"access": "public", "ctype": "Struct", "type": "int"}),
		tag("Field2", 4, "w", F{"access": "public", "ctype": "Struct", "type": "int"}),
		tag("field3", 5, "w", F{"access": "private", "ctype": "Struct", "type": "string"}),
		tag("field4", 6, "w", F{"access": "private", "ctype": "Struct", "type": "*bool"}),
		tag("Struct", 3, "t", F{"access": "public", "type": "struct"}),
		tag("Struct", 20, "e", F{"access": "public", "ctype": "TestEmbed", "type": "Struct"}),
		tag("*io.Writer", 21, "e", F{"access": "public", "ctype": "TestEmbed", "type": "*io.Writer"}),
		tag("TestEmbed", 19, "t", F{"access": "public", "type": "struct"}),
		tag("Struct2", 27, "t", F{"access": "public", "type": "struct"}),
		tag("Connection", 36, "t", F{"access": "public", "type": "struct"}),
		tag("NewStruct", 9, "f", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "*Struct"}),
		tag("F1", 13, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "[]bool, [2]*string"}),
		tag("F2", 16, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "bool"}),
		tag("NewTestEmbed", 24, "f", F{"access": "public", "ctype": "TestEmbed", "signature": "()", "type": "TestEmbed"}),
		tag("NewStruct2", 30, "f", F{"access": "public", "ctype": "Struct2", "signature": "()", "type": "*Struct2, error"}),
		tag("Dial", 33, "f", F{"access": "public", "ctype": "Connection", "signature": "()", "type": "*Connection, error"}),
		tag("Dial2", 39, "f", F{"access": "public", "ctype": "Connection", "signature": "()", "type": "*Connection, *Struct2"}),
		tag("Dial3", 42, "f", F{"access": "public", "signature": "()", "type": "*Connection, *Connection"}),
	}},
	{filename: "testdata/struct.go", withExtraSymbols: true, tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Field1", 4, "w", F{"access": "public", "ctype": "Struct", "type": "int"}),
		tag("Field2", 4, "w", F{"access": "public", "ctype": "Struct", "type": "int"}),
		tag("field3", 5, "w", F{"access": "private", "ctype": "Struct", "type": "string"}),
		tag("field4", 6, "w", F{"access": "private", "ctype": "Struct", "type": "*bool"}),
		tag("Struct", 3, "t", F{"access": "public", "type": "struct"}),
		tag("Test.Struct", 3, "t", F{"access": "public", "type": "struct"}),
		tag("Struct", 20, "e", F{"access": "public", "ctype": "TestEmbed", "type": "Struct"}),
		tag("*io.Writer", 21, "e", F{"access": "public", "ctype": "TestEmbed", "type": "*io.Writer"}),
		tag("TestEmbed", 19, "t", F{"access": "public", "type": "struct"}),
		tag("Test.TestEmbed", 19, "t", F{"access": "public", "type": "struct"}),
		tag("Struct2", 27, "t", F{"access": "public", "type": "struct"}),
		tag("Test.Struct2", 27, "t", F{"access": "public", "type": "struct"}),
		tag("Connection", 36, "t", F{"access": "public", "type": "struct"}),
		tag("Test.Connection", 36, "t", F{"access": "public", "type": "struct"}),
		tag("NewStruct", 9, "f", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "*Struct"}),
		tag("Test.NewStruct", 9, "f", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "*Struct"}),
		tag("F1", 13, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "[]bool, [2]*string"}),
		tag("Struct.F1", 13, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "[]bool, [2]*string"}),
		tag("Test.F1", 13, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "[]bool, [2]*string"}),
		tag("Test.Struct.F1", 13, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "[]bool, [2]*string"}),
		tag("F2", 16, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "bool"}),
		tag("Struct.F2", 16, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "bool"}),
		tag("Test.Struct.F2", 16, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "bool"}),
		tag("Test.F2", 16, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "bool"}),
		tag("NewTestEmbed", 24, "f", F{"access": "public", "ctype": "TestEmbed", "signature": "()", "type": "TestEmbed"}),
		tag("Test.NewTestEmbed", 24, "f", F{"access": "public", "ctype": "TestEmbed", "signature": "()", "type": "TestEmbed"}),
		tag("NewStruct2", 30, "f", F{"access": "public", "ctype": "Struct2", "signature": "()", "type": "*Struct2, error"}),
		tag("Test.NewStruct2", 30, "f", F{"access": "public", "ctype": "Struct2", "signature": "()", "type": "*Struct2, error"}),
		tag("Dial", 33, "f", F{"access": "public", "ctype": "Connection", "signature": "()", "type": "*Connection, error"}),
		tag("Test.Dial", 33, "f", F{"access": "public", "ctype": "Connection", "signature": "()", "type": "*Connection, error"}),
		tag("Dial2", 39, "f", F{"access": "public", "ctype": "Connection", "signature": "()", "type": "*Connection, *Struct2"}),
		tag("Test.Dial2", 39, "f", F{"access": "public", "ctype": "Connection", "signature": "()", "type": "*Connection, *Struct2"}),
		tag("Dial3", 42, "f", F{"access": "public", "signature": "()", "type": "*Connection, *Connection"}),
		tag("Test.Dial3", 42, "f", F{"access": "public", "signature": "()", "type": "*Connection, *Connection"}),
	}},
	{filename: "testdata/type.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("testType", 3, "t", F{"access": "private", "type": "int"}),
		tag("testArrayType", 4, "t", F{"access": "private", "type": "[4]int"}),
		tag("testSliceType", 5, "t", F{"access": "private", "type": "[]int"}),
		tag("testPointerType", 6, "t", F{"access": "private", "type": "*string"}),
		tag("testFuncType1", 7, "t", F{"access": "private", "type": "func()"}),
		tag("testFuncType2", 8, "t", F{"access": "private", "type": "func(int) string"}),
		tag("testMapType", 9, "t", F{"access": "private", "type": "map[string]bool"}),
		tag("testChanType", 10, "t", F{"access": "private", "type": "chan bool"}),
	}},
	{filename: "testdata/type.go", withExtraSymbols: true, tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("testType", 3, "t", F{"access": "private", "type": "int"}),
		tag("testArrayType", 4, "t", F{"access": "private", "type": "[4]int"}),
		tag("testSliceType", 5, "t", F{"access": "private", "type": "[]int"}),
		tag("testPointerType", 6, "t", F{"access": "private", "type": "*string"}),
		tag("testFuncType1", 7, "t", F{"access": "private", "type": "func()"}),
		tag("testFuncType2", 8, "t", F{"access": "private", "type": "func(int) string"}),
		tag("testMapType", 9, "t", F{"access": "private", "type": "map[string]bool"}),
		tag("testChanType", 10, "t", F{"access": "private", "type": "chan bool"}),
		tag("Test.testType", 3, "t", F{"access": "private", "type": "int"}),
		tag("Test.testArrayType", 4, "t", F{"access": "private", "type": "[4]int"}),
		tag("Test.testSliceType", 5, "t", F{"access": "private", "type": "[]int"}),
		tag("Test.testPointerType", 6, "t", F{"access": "private", "type": "*string"}),
		tag("Test.testFuncType1", 7, "t", F{"access": "private", "type": "func()"}),
		tag("Test.testFuncType2", 8, "t", F{"access": "private", "type": "func(int) string"}),
		tag("Test.testMapType", 9, "t", F{"access": "private", "type": "map[string]bool"}),
		tag("Test.testChanType", 10, "t", F{"access": "private", "type": "chan bool"}),
	}},
	{filename: "testdata/var.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("variable1", 3, "v", F{"access": "private", "type": "int"}),
		tag("variable2", 4, "v", F{"access": "private", "type": "string"}),
		tag("A", 7, "v", F{"access": "public"}),
		tag("B", 8, "v", F{"access": "public"}),
		tag("C", 8, "v", F{"access": "public"}),
		tag("D", 9, "v", F{"access": "public"}),
	}},
	{filename: "testdata/var.go", withExtraSymbols: true, tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("variable1", 3, "v", F{"access": "private", "type": "int"}),
		tag("variable2", 4, "v", F{"access": "private", "type": "string"}),
		tag("A", 7, "v", F{"access": "public"}),
		tag("B", 8, "v", F{"access": "public"}),
		tag("C", 8, "v", F{"access": "public"}),
		tag("D", 9, "v", F{"access": "public"}),
		tag("Test.variable1", 3, "v", F{"access": "private", "type": "int"}),
		tag("Test.variable2", 4, "v", F{"access": "private", "type": "string"}),
		tag("Test.A", 7, "v", F{"access": "public"}),
		tag("Test.B", 8, "v", F{"access": "public"}),
		tag("Test.C", 8, "v", F{"access": "public"}),
		tag("Test.D", 9, "v", F{"access": "public"}),
	}},
	{filename: "testdata/simple.go", relative: true, basepath: "dir", tags: []Tag{
		{Name: "main", File: "../testdata/simple.go", Address: "1", Type: "p", Fields: F{"line": "1"}},
	}},
	{filename: "testdata/simple.go", withExtraSymbols: true, relative: true, basepath: "dir", tags: []Tag{
		{Name: "main", File: "../testdata/simple.go", Address: "1", Type: "p", Fields: F{"line": "1"}},
	}},
	{filename: "testdata/range.go", minversion: 4, tags: []Tag{
		tag("main", 1, "p", F{}),
		tag("fmt", 3, "i", F{}),
		tag("main", 5, "f", F{"access": "private", "signature": "()"}),
	}},
	{filename: "testdata/range.go", withExtraSymbols: true, minversion: 4, tags: []Tag{
		tag("main", 1, "p", F{}),
		tag("fmt", 3, "i", F{}),
		tag("main", 5, "f", F{"access": "private", "signature": "()"}),
		tag("main.main", 5, "f", F{"access": "private", "signature": "()"}),
	}},
}

func TestParse(t *testing.T) {
	for _, testCase := range testCases {
		if testCase.minversion > 0 && extractVersionCode(runtime.Version()) < testCase.minversion {
			t.Skipf("[%s] skipping test. Version is %s, but test requires at least go1.%d", testCase.filename, runtime.Version(), testCase.minversion)
			continue
		}

		basepath, err := filepath.Abs(testCase.basepath)
		if err != nil {
			t.Errorf("[%s] could not determine base path: %s\n", testCase.filename, err)
			continue
		}

		var extra FieldSet
		if testCase.withExtraSymbols {
			extra = FieldSet{ExtraTags: true}
		}

		tags, err := Parse(testCase.filename, testCase.relative, basepath, extra)
		if err != nil {
			t.Errorf("[%s] Parse error: %s", testCase.filename, err)
			continue
		}

		sort.Sort(TagSlice(tags))
		sort.Sort(TagSlice(testCase.tags))

		if len(tags) != len(testCase.tags) {
			t.Errorf("[%s] len(tags) == %d, want %d", testCase.filename, len(tags), len(testCase.tags))
			continue
		}

		for i, tag := range testCase.tags {
			if len(tag.File) == 0 {
				tag.File = testCase.filename
			}
			if tags[i].String() != tag.String() {
				t.Errorf("[%s] tag(%d)\n  is:%s\nwant:%s", testCase.filename, i, tags[i].String(), tag.String())
			}
		}
	}
}

func tag(n string, l int, t TagType, fields F) (tag Tag) {
	tag = Tag{
		Name:    n,
		File:    "",
		Address: strconv.Itoa(l),
		Type:    t,
		Fields:  fields,
	}

	tag.Fields["line"] = tag.Address

	return
}

func extractVersionCode(version string) int {
	matches := goVersionRegexp.FindAllStringSubmatch(version, -1)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(matches[0][1])
	return n
}

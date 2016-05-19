package main

import (
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

type F map[TagField]string

var testCases = []struct {
	filename   string
	relative   bool
	basepath   string
	minversion string
	tags       []Tag
}{
	{filename: "tests/const.go-src", tags: []Tag{
		tag("Test", 1, 9, "p", F{}),
		tag("Constant", 3, 7, "c", F{"access": "public", "type": "string"}),
		tag("OtherConst", 4, 7, "c", F{"access": "public"}),
		tag("A", 7, 2, "c", F{"access": "public"}),
		tag("B", 8, 2, "c", F{"access": "public"}),
		tag("C", 8, 5, "c", F{"access": "public"}),
		tag("D", 9, 5, "c", F{"access": "public"}),
	}},
	{filename: "tests/func.go-src", tags: []Tag{
		tag("Test", 1, 9, "p", F{}),
		tag("Function1", 3, 1, "f", F{"access": "public", "signature": "()", "type": "string"}),
		tag("function2", 6, 1, "f", F{"access": "private", "signature": "(p1, p2 int, p3 *string)"}),
		tag("function3", 9, 1, "f", F{"access": "private", "signature": "()", "type": "bool"}),
		tag("function4", 12, 1, "f", F{"access": "private", "signature": "(p interface{})", "type": "interface{}"}),
		tag("function5", 15, 1, "f", F{"access": "private", "signature": "()", "type": "string, string, error"}),
	}},
	{filename: "tests/import.go-src", tags: []Tag{
		tag("Test", 1, 9, "p", F{}),
		tag("fmt", 3, 8, "i", F{}),
		tag("go/ast", 6, 2, "i", F{}),
		tag("go/parser", 7, 2, "i", F{}),
	}},
	{filename: "tests/interface.go-src", tags: []Tag{
		tag("Test", 1, 9, "p", F{}),
		tag("InterfaceMethod", 4, 2, "m", F{"access": "public", "signature": "(int)", "ntype": "Interface", "type": "string"}),
		tag("OtherMethod", 5, 2, "m", F{"access": "public", "signature": "()", "ntype": "Interface"}),
		tag("io.Reader", 6, 2, "e", F{"access": "public", "ntype": "Interface"}),
		tag("Interface", 3, 6, "n", F{"access": "public", "type": "interface"}),
	}},
	{filename: "tests/struct.go-src", tags: []Tag{
		tag("Test", 1, 9, "p", F{}),
		tag("Field1", 4, 2, "w", F{"access": "public", "ctype": "Struct", "type": "int"}),
		tag("Field2", 4, 10, "w", F{"access": "public", "ctype": "Struct", "type": "int"}),
		tag("field3", 5, 2, "w", F{"access": "private", "ctype": "Struct", "type": "string"}),
		tag("field4", 6, 2, "w", F{"access": "private", "ctype": "Struct", "type": "*bool"}),
		tag("Struct", 3, 6, "t", F{"access": "public", "type": "struct"}),
		tag("Struct", 20, 2, "e", F{"access": "public", "ctype": "TestEmbed", "type": "Struct"}),
		tag("*io.Writer", 21, 2, "e", F{"access": "public", "ctype": "TestEmbed", "type": "*io.Writer"}),
		tag("TestEmbed", 19, 6, "t", F{"access": "public", "type": "struct"}),
		tag("Struct2", 27, 6, "t", F{"access": "public", "type": "struct"}),
		tag("Connection", 36, 6, "t", F{"access": "public", "type": "struct"}),
		tag("NewStruct", 9, 1, "f", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "*Struct"}),
		tag("F1", 13, 1, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "[]bool, [2]*string", "receiver": "s"}),
		tag("F2", 16, 1, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "bool"}),
		tag("NewTestEmbed", 24, 1, "f", F{"access": "public", "ctype": "TestEmbed", "signature": "()", "type": "TestEmbed"}),
		tag("NewStruct2", 30, 1, "f", F{"access": "public", "ctype": "Struct2", "signature": "()", "type": "*Struct2, error"}),
		tag("Dial", 33, 1, "f", F{"access": "public", "ctype": "Connection", "signature": "()", "type": "*Connection, error"}),
		tag("Dial2", 39, 1, "f", F{"access": "public", "ctype": "Connection", "signature": "()", "type": "*Connection, *Struct2"}),
		tag("Dial3", 42, 1, "f", F{"access": "public", "signature": "()", "type": "*Connection, *Connection"}),
	}},
	{filename: "tests/type.go-src", tags: []Tag{
		tag("Test", 1, 9, "p", F{}),
		tag("testType", 3, 6, "t", F{"access": "private", "type": "int"}),
		tag("testArrayType", 4, 6, "t", F{"access": "private", "type": "[4]int"}),
		tag("testSliceType", 5, 6, "t", F{"access": "private", "type": "[]int"}),
		tag("testPointerType", 6, 6, "t", F{"access": "private", "type": "*string"}),
		tag("testFuncType1", 7, 6, "t", F{"access": "private", "type": "func()"}),
		tag("testFuncType2", 8, 6, "t", F{"access": "private", "type": "func(int) string"}),
		tag("testMapType", 9, 6, "t", F{"access": "private", "type": "map[string]bool"}),
		tag("testChanType", 10, 6, "t", F{"access": "private", "type": "chan bool"}),
	}},
	{filename: "tests/var.go-src", tags: []Tag{
		tag("Test", 1, 9, "p", F{}),
		tag("variable1", 3, 5, "v", F{"access": "private", "type": "int"}),
		tag("variable2", 4, 5, "v", F{"access": "private", "type": "string"}),
		tag("A", 7, 2, "v", F{"access": "public"}),
		tag("B", 8, 2, "v", F{"access": "public"}),
		tag("C", 8, 5, "v", F{"access": "public"}),
		tag("D", 9, 5, "v", F{"access": "public"}),
	}},
	{filename: "tests/simple.go-src", relative: true, basepath: "dir", tags: []Tag{
		Tag{Name: "main", File: "../tests/simple.go-src", Address: "1", Type: "p", Fields: F{"line": "1", "column": "9"}},
	}},
	{filename: "tests/range.go-src", minversion: "go1.4", tags: []Tag{
		tag("main", 1, 9, "p", F{}),
		tag("fmt", 3, 8, "i", F{}),
		tag("main", 5, 1, "f", F{"access": "private", "signature": "()"}),
	}},
	{filename: "tests/method.go-src", tags: []Tag{
		tag("main", 1, 9, "p", F{}),
		tag("Struct", 3, 6, "t", F{"access": "public", "type": "struct"}),
		tag("methodA", 8, 1, "m", F{"access": "private", "receiver": "s", "ctype": "Struct", "signature": "()"}),
		tag("s", 10, 2, "x", F{}),
		tag("s", 11, 17, "x", F{}),
		tag("s", 16, 12, "x", F{}),
		tag("s", 17, 8, "x", F{}),
	}},
}

func TestParse(t *testing.T) {
	for _, testCase := range testCases {
		if testCase.minversion != "" && runtime.Version() < testCase.minversion {
			t.Skipf("[%s] skipping test. Version is %s, but test requires %s", testCase.filename, runtime.Version(), testCase.minversion)
			continue
		}

		basepath, err := filepath.Abs(testCase.basepath)
		if err != nil {
			t.Errorf("[%s] could not determine base path: %s\n", testCase.filename, err)
			continue
		}

		tags, err := Parse(testCase.filename, testCase.relative, basepath)
		if err != nil {
			t.Errorf("[%s] Parse error: %s", testCase.filename, err)
			continue
		}

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

func tag(n string, l, c int, t TagType, fields F) (tag Tag) {
	tag = Tag{
		Name:    n,
		File:    "",
		Address: strconv.Itoa(l),
		Type:    t,
		Fields:  fields,
	}

	tag.Fields["line"] = tag.Address
	tag.Fields["column"] = strconv.Itoa(c)

	return
}

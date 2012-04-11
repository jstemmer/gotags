package main

import (
	"strconv"
	"testing"
)

type F map[string]string

var testCases = []struct {
	filename string
	tags     []Tag
}{
	{filename: "tests/const.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Constant", 3, "c", F{"access": "public", "type": "string"}),
		tag("OtherConst", 4, "c", F{"access": "public"}),
		tag("A", 7, "c", F{"access": "public"}),
		tag("B", 8, "c", F{"access": "public"}),
		tag("C", 8, "c", F{"access": "public"}),
		tag("D", 9, "c", F{"access": "public"}),
	}},
	{filename: "tests/func.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Function1", 3, "f", F{"access": "public", "signature": "()", "type": "string"}),
		tag("function2", 6, "f", F{"access": "private", "signature": "(p1, p2 int, p3 *string)"}),
		tag("function3", 9, "f", F{"access": "private", "signature": "()", "type": "bool"}),
		tag("function4", 12, "f", F{"access": "private", "signature": "(p interface{})", "type": "interface{}"}),
	}},
	{filename: "tests/import.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("fmt", 3, "i", F{}),
		tag("go/ast", 6, "i", F{}),
		tag("go/parser", 7, "i", F{}),
	}},
	{filename: "tests/interface.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("InterfaceMethod", 4, "m", F{"access": "public", "signature": "(int)", "ntype": "Interface", "type": "string"}),
		tag("OtherMethod", 5, "m", F{"access": "public", "signature": "()", "ntype": "Interface"}),
		tag("io.Reader", 6, "e", F{"access": "public", "ntype": "Interface"}),
		tag("Interface", 3, "n", F{"access": "public", "type": "interface"}),
	}},
	{filename: "tests/struct.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("Field1", 4, "w", F{"access": "public", "ctype": "Struct", "type": "int"}),
		tag("Field2", 4, "w", F{"access": "public", "ctype": "Struct", "type": "int"}),
		tag("field3", 5, "w", F{"access": "private", "ctype": "Struct", "type": "string"}),
		tag("field4", 6, "w", F{"access": "private", "ctype": "Struct", "type": "*bool"}),
		tag("Struct", 3, "t", F{"access": "public", "type": "struct"}),
		tag("NewStruct", 9, "r", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "*Struct"}),
		tag("F1", 13, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "[]bool, [2]*string"}),
		tag("F2", 16, "m", F{"access": "public", "ctype": "Struct", "signature": "()", "type": "bool"}),
		tag("Struct", 20, "e", F{"access": "public", "ctype": "TestEmbed", "type": "Struct"}),
		tag("*io.Writer", 21, "e", F{"access": "public", "ctype": "TestEmbed", "type": "*io.Writer"}),
		tag("TestEmbed", 19, "t", F{"access": "public", "type": "struct"}),
		tag("NewTestEmbed", 24, "r", F{"access": "public", "ctype": "TestEmbed", "signature": "()", "type": "TestEmbed"}),
		tag("Struct2", 27, "t", F{"access": "public", "type": "struct"}),
		tag("NewStruct2", 30, "r", F{"access": "public", "ctype": "Struct2", "signature": "()", "type": "*Struct2, error"}),
	}},
	{filename: "tests/type.go", tags: []Tag{
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
	{filename: "tests/var.go", tags: []Tag{
		tag("Test", 1, "p", F{}),
		tag("variable1", 3, "v", F{"access": "private", "type": "int"}),
		tag("variable2", 4, "v", F{"access": "private", "type": "string"}),
		tag("A", 7, "v", F{"access": "public"}),
		tag("B", 8, "v", F{"access": "public"}),
		tag("C", 8, "v", F{"access": "public"}),
		tag("D", 9, "v", F{"access": "public"}),
	}},
}

func TestParse(t *testing.T) {
	for _, testCase := range testCases {
		tags, err := Parse(testCase.filename)
		if err != nil {
			t.Errorf("[%s] Parse error: %s", testCase.filename, err)
			continue
		}

		if len(tags) != len(testCase.tags) {
			t.Errorf("[%s] len(tags) == %d, want %d", testCase.filename, len(tags), len(testCase.tags))
			continue
		}

		for i, tag := range testCase.tags {
			tag.File = testCase.filename
			if tags[i].String() != tag.String() {
				t.Errorf("[%s] tag(%d)\n  is:%s\nwant:%s", testCase.filename, i, tags[i].String(), tag.String())
			}
		}
	}
}

func tag(n string, l int, t string, fields F) (tag Tag) {
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

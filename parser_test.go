package main

import (
	"testing"
)

var expectedTags = []Tag{
	Tag{Name: "TestPackage", File: "tests/input.go", Address: "1", Type: "p", Fields: map[string]string{"line": "1"}},
	Tag{Name: "go/ast", File: "tests/input.go", Address: "4", Type: "i", Fields: map[string]string{"line": "4"}},
	Tag{Name: "variable", File: "tests/input.go", Address: "7", Type: "v", Fields: map[string]string{"line": "7", "access": "private", "type": "int"}},
	Tag{Name: "Constant", File: "tests/input.go", Address: "9", Type: "c", Fields: map[string]string{"line": "9", "access": "public", "type": "string"}},
	Tag{Name: "OtherConst", File: "tests/input.go", Address: "10", Type: "c", Fields: map[string]string{"line": "10", "access": "public"}},
	Tag{Name: "Function1", File: "tests/input.go", Address: "12", Type: "f", Fields: map[string]string{"line": "12", "access": "public", "signature": "()", "type": "string"}},
	Tag{Name: "function2", File: "tests/input.go", Address: "15", Type: "f", Fields: map[string]string{"line": "15", "access": "private", "signature": "(p1, p2 int, p3 *string)"}},
	Tag{Name: "Struct", File: "tests/input.go", Address: "18", Type: "t", Fields: map[string]string{"line": "18", "access": "public", "type": "struct"}},
	Tag{Name: "Field1", File: "tests/input.go", Address: "19", Type: "w", Fields: map[string]string{"line": "19", "access": "public", "ctype": "Struct", "type": "int"}},
	Tag{Name: "field2", File: "tests/input.go", Address: "20", Type: "w", Fields: map[string]string{"line": "20", "access": "private", "ctype": "Struct", "type": "string"}},
	Tag{Name: "field3", File: "tests/input.go", Address: "21", Type: "w", Fields: map[string]string{"line": "21", "access": "private", "ctype": "Struct", "type": "*bool"}},
	Tag{Name: "myInt", File: "tests/input.go", Address: "24", Type: "t", Fields: map[string]string{"line": "24", "access": "private", "type": "int"}},
	Tag{Name: "F1", File: "tests/input.go", Address: "26", Type: "f", Fields: map[string]string{"line": "26", "access": "public", "signature": "()", "ctype": "myInt", "type": "[]bool, [2]*string"}},
	Tag{Name: "TestEmbed", File: "tests/input.go", Address: "29", Type: "t", Fields: map[string]string{"line": "29", "access": "public", "type": "struct"}},
	Tag{Name: "Struct", File: "tests/input.go", Address: "30", Type: "w", Fields: map[string]string{"line": "30", "access": "public", "ctype": "TestEmbed", "type": "Struct"}},
}

func TestParse(t *testing.T) {
	tags, err := Parse("tests/input.go")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if len(tags) != len(expectedTags) {
		t.Fatalf("len(tags) == %d, want %d", len(tags), len(expectedTags))
	}

	for i, exp := range expectedTags {
		if tags[i].String() != exp.String() {
			t.Errorf("tag(%d)\n  is:%s\nwant:%s", i, tags[i].String(), exp.String())
		}
	}
}

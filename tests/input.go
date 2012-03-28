package TestPackage

import (
	"go/ast"
)

var variable int

const constant = "const"

func Function1() string {
}

func function2(p1, p2 int, p3 *string) {
}

type Struct struct {
	Field1 int
	field2 string
	field3 *bool
}

type MyInt int

func (m MyInt) F1() {
}

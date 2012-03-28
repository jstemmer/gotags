package TestPackage

import (
	"go/ast"
)

var variable int

const Constant = "const"

func Function1() string {
}

func function2(p1, p2 int, p3 *string) {
}

type Struct struct {
	Field1 int
	field2 string
	field3 *bool
}

type myInt int

func (m myInt) F1() {
}

type TestEmbed struct {
	Struct
}

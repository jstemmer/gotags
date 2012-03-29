package TestPackage

import (
	"go/ast"
)

var variable int

const Constant string = "const"
const OtherConst = "const"

func Function1() string {
}

func function2(p1, p2 int, p3 *string) {
}

type Struct struct {
	Field1, Field2 int
	field3         string
	field4         *bool
}

func NewStruct() *Struct {
	return &Struct{}
}

type myInt int

func (m myInt) F1() ([]bool, [2]*string) {
}

type TestEmbed struct {
	Struct
	*io.Writer
}

type Interface interface {
	InterfaceMethod(int) string
	OtherMethod()
	io.Reader
}

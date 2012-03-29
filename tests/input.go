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

type testArrayType [4]int
type testSliceType []int
type testPointerType *string
type testFuncType1 func()
type testFuncType2 func(int) string
type testMapType map[string]bool
type testChanType chan bool

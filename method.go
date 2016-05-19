package main

import "go/ast"

// methodNodeVisitor implements interface ast.Visitor and used for walking of
// []ast.Node for obtaining receivers.
type methodNodeVisitor struct {
	receiver *ast.Object
	parser   *tagParser
}

// Visit investigates specified ast.Node for usage of specified method receiver.
func (visitor *methodNodeVisitor) Visit(node ast.Node) ast.Visitor {
	switch stmt := node.(type) {
	case *ast.Ident:
		if stmt.Obj == visitor.receiver {
			tag := visitor.parser.createTag(stmt.Name, stmt.Pos(), Receiver)
			visitor.parser.tags = append(visitor.parser.tags, tag)
		}
	}
	return visitor
}

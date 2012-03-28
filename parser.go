package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type tagParser struct {
	fset *token.FileSet
	tags []Tag
}

func Parse(filename string) ([]Tag, error) {
	parser := &tagParser{
		fset: token.NewFileSet(),
		tags: make([]Tag, 0),
	}

	err := parser.parseFile(filename)
	if err != nil {
		return nil, err
	}

	return parser.tags, nil
}

func (p *tagParser) parseFile(filename string) error {
	f, err := parser.ParseFile(p.fset, filename, nil, 0)
	if err != nil {
		return err
	}

	// package
	p.parsePackage(f)

	// imports
	p.parseImports(f)

	// declarations
	p.parseDeclarations(f)

	return nil
}

func (p *tagParser) parsePackage(f *ast.File) {
	p.tags = append(p.tags, p.createTag(f.Name.Name, f.Name.Pos(), "p"))
}

func (p *tagParser) parseImports(f *ast.File) {
	for _, im := range f.Imports {
		name := strings.Trim(im.Path.Value, "\"")
		p.tags = append(p.tags, p.createTag(name, im.Path.Pos(), "i"))
	}
}

func (p *tagParser) parseDeclarations(f *ast.File) {
	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			p.parseFunction(decl)
		case *ast.GenDecl:
			for _, s := range decl.Specs {
				switch ts := s.(type) {
				case *ast.TypeSpec:
					tag := p.createTag(ts.Name.Name, ts.Pos(), "t")
					if ast.IsExported(tag.Name) {
						tag.Fields["access"] = "public"
					} else {
						tag.Fields["access"] = "private"
					}
					p.tags = append(p.tags, tag)
				case *ast.ValueSpec:
					if len(ts.Names) > 0 {
						tag := p.createTag(ts.Names[0].Name, ts.Pos(), "v")

						switch ts.Names[0].Obj.Kind {
						case ast.Var:
							tag.Type = "v"
						case ast.Con:
							tag.Type = "c"
						}

						p.tags = append(p.tags, tag)
					}
				}
			}
		}
	}
}

func (p *tagParser) parseFunction(f *ast.FuncDecl) {
	tag := p.createTag(f.Name.Name, f.Pos(), "f")

	// access
	if ast.IsExported(tag.Name) {
		tag.Fields["access"] = "public"
	} else {
		tag.Fields["access"] = "private"
	}

	// signature
	var sig bytes.Buffer
	sig.WriteByte('(')
	for i, param := range f.Type.Params.List {
		// parameter names
		for j, n := range param.Names {
			sig.WriteString(n.Name)
			if j < len(param.Names)-1 {
				sig.WriteString(", ")
			}
		}

		// parameter type
		sig.WriteByte(' ')
		sig.WriteString(getType(param.Type, true))

		if i < len(f.Type.Params.List)-1 {
			sig.WriteString(", ")
		}
	}
	sig.WriteByte(')')
	tag.Fields["signature"] = sig.String()

	// receiver
	if f.Recv != nil && len(f.Recv.List) > 0 {
		tag.Fields["type"] = getType(f.Recv.List[0].Type, false)
	}

	p.tags = append(p.tags, tag)
}

func (p *tagParser) createTag(name string, pos token.Pos, tagtype string) Tag {
	return NewTag(name, p.fset.File(pos).Name(), p.fset.Position(pos).Line, tagtype)
}

func getType(node ast.Node, star bool) (paramType string) {
	switch t := node.(type) {
	case *ast.Ident:
		paramType = t.Name
	case *ast.StarExpr:
		if star {
			paramType = "*" + getType(t.X, star)
		} else {
			paramType = getType(t.X, star)
		}
	case *ast.SelectorExpr:
		paramType = getType(t.X, star) + "." + getType(t.Sel, star)
	}
	return
}

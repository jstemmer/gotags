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

func PrintTree(filename string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return err
	}

	ast.Print(fset, f)
	return nil
}

func Parse(filename string) ([]Tag, error) {
	p := &tagParser{
		fset: token.NewFileSet(),
		tags: make([]Tag, 0),
	}

	err := p.parseFile(filename)
	if err != nil {
		return nil, err
	}

	return p.tags, nil
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
					p.parseTypeDeclaration(ts)
				case *ast.ValueSpec:
					p.parseValueDeclaration(ts)
				}
			}
		}
	}
}

func (p *tagParser) parseFunction(f *ast.FuncDecl) {
	tag := p.createTag(f.Name.Name, f.Pos(), "f")

	// access
	tag.Fields["access"] = getAccess(tag.Name)

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

func (p *tagParser) parseTypeDeclaration(ts *ast.TypeSpec) {
	tag := p.createTag(ts.Name.Name, ts.Pos(), "t")

	tag.Fields["access"] = getAccess(tag.Name)

	p.tags = append(p.tags, tag)

	if s, ok := ts.Type.(*ast.StructType); ok {
		p.parseStructFields(tag.Name, s)
	}
}

func (p *tagParser) parseValueDeclaration(v *ast.ValueSpec) {
	tag := p.createTag(v.Names[0].Name, v.Pos(), "v")

	tag.Fields["access"] = getAccess(tag.Name)

	switch v.Names[0].Obj.Kind {
	case ast.Var:
		tag.Type = "v"
	case ast.Con:
		tag.Type = "c"
	}

	p.tags = append(p.tags, tag)
}

func (p *tagParser) parseStructFields(name string, s *ast.StructType) {
	for _, f := range s.Fields.List {
		var tag Tag
		if len(f.Names) > 0 {
			tag = p.createTag(f.Names[0].Name, f.Names[0].Pos(), "w")
		} else if t, ok := f.Type.(*ast.Ident); ok {
			tag = p.createTag(t.Name, t.Pos(), "w")
		}

		tag.Fields["access"] = getAccess(tag.Name)
		tag.Fields["type"] = name

		p.tags = append(p.tags, tag)
	}
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

func getAccess(name string) (access string) {
	if ast.IsExported(name) {
		access = "public"
	} else {
		access = "private"
	}
	return
}

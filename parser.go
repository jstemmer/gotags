package main

import (
	"fmt"
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
		tags: []Tag{},
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

	tag.Fields["access"] = getAccess(tag.Name)
	tag.Fields["signature"] = fmt.Sprintf("(%s)", getTypes(f.Type.Params.List, true))

	if f.Type.Results != nil {
		tag.Fields["type"] = getTypes(f.Type.Results.List, false)
	}

	// receiver
	if f.Recv != nil && len(f.Recv.List) > 0 {
		tag.Fields["ctype"] = getType(f.Recv.List[0].Type, false)
		tag.Type = "m"
	}

	// check if this is a constructor, in that case it belongs to that type
	if strings.HasPrefix(tag.Name, "New") && len(tag.Fields["type"]) > 0 {
		if tag.Name[3:] == tag.Fields["type"] {
			tag.Fields["ctype"] = tag.Fields["type"]
			tag.Type = "r"
		} else if tag.Fields["type"][0] == '*' && tag.Name[3:] == tag.Fields["type"][1:] {
			tag.Fields["ctype"] = tag.Fields["type"][1:]
			tag.Type = "r"
		}
	}

	p.tags = append(p.tags, tag)
}

func (p *tagParser) parseTypeDeclaration(ts *ast.TypeSpec) {
	tag := p.createTag(ts.Name.Name, ts.Pos(), "t")

	tag.Fields["access"] = getAccess(tag.Name)

	switch s := ts.Type.(type) {
	case *ast.StructType:
		tag.Fields["type"] = "struct"
		p.parseStructFields(tag.Name, s)
	case *ast.InterfaceType:
		tag.Fields["type"] = "interface"
		tag.Type = "n"
		p.parseInterfaceMethods(tag.Name, s)
	default:
		tag.Fields["type"] = getType(ts.Type, true)
	}

	p.tags = append(p.tags, tag)
}

func (p *tagParser) parseValueDeclaration(v *ast.ValueSpec) {
	for _, d := range v.Names {
		if d.Name == "_" {
			continue
		}

		tag := p.createTag(d.Name, d.Pos(), "v")
		tag.Fields["access"] = getAccess(tag.Name)

		if v.Type != nil {
			tag.Fields["type"] = getType(v.Type, true)
		}

		switch d.Obj.Kind {
		case ast.Var:
			tag.Type = "v"
		case ast.Con:
			tag.Type = "c"
		}
		p.tags = append(p.tags, tag)
	}
}

func (p *tagParser) parseStructFields(name string, s *ast.StructType) {
	for _, f := range s.Fields.List {
		var tag Tag
		if len(f.Names) > 0 {
			for _, n := range f.Names {
				tag = p.createTag(n.Name, n.Pos(), "w")
				tag.Fields["access"] = getAccess(tag.Name)
				tag.Fields["ctype"] = name
				tag.Fields["type"] = getType(f.Type, true)
				p.tags = append(p.tags, tag)
			}
		} else {
			// embedded field
			tag = p.createTag(getType(f.Type, true), f.Pos(), "e")
			tag.Fields["access"] = getAccess(tag.Name)
			tag.Fields["ctype"] = name
			tag.Fields["type"] = getType(f.Type, true)
			p.tags = append(p.tags, tag)
		}
	}
}

func (p *tagParser) parseInterfaceMethods(name string, s *ast.InterfaceType) {
	for _, f := range s.Methods.List {
		var tag Tag
		if len(f.Names) > 0 {
			tag = p.createTag(f.Names[0].Name, f.Names[0].Pos(), "m")
		} else {
			// embedded interface
			tag = p.createTag(getType(f.Type, true), f.Pos(), "e")
		}

		tag.Fields["access"] = getAccess(tag.Name)

		if t, ok := f.Type.(*ast.FuncType); ok {
			tag.Fields["signature"] = fmt.Sprintf("(%s)", getTypes(t.Params.List, true))
			if t.Results != nil {
				tag.Fields["type"] = getTypes(t.Results.List, false)
			}
		}

		tag.Fields["ntype"] = name

		p.tags = append(p.tags, tag)
	}
}

func (p *tagParser) createTag(name string, pos token.Pos, tagtype string) Tag {
	return NewTag(name, p.fset.File(pos).Name(), p.fset.Position(pos).Line, tagtype)
}

func getTypes(fields []*ast.Field, includeNames bool) string {
	types := make([]string, len(fields))
	for i, param := range fields {
		if includeNames && len(param.Names) > 0 {
			// parameter names
			names := make([]string, len(param.Names))
			for j, n := range param.Names {
				names[j] = n.Name
			}
			types[i] = fmt.Sprintf("%s %s", strings.Join(names, ", "), getType(param.Type, true))
		} else {
			types[i] = getType(param.Type, true)
		}
	}

	return strings.Join(types, ", ")
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
	case *ast.ArrayType:
		if l, ok := t.Len.(*ast.BasicLit); ok {
			paramType = fmt.Sprintf("[%s]%s", l.Value, getType(t.Elt, star))
		} else {
			paramType = "[]" + getType(t.Elt, star)
		}
	case *ast.FuncType:
		var fparams, fresult string
		if t.Params != nil {
			fparams = getTypes(t.Params.List, true)
		}

		if t.Results != nil {
			fresult = getTypes(t.Results.List, false)
		}

		if len(fresult) > 0 {
			paramType = fmt.Sprintf("func(%s) %s", fparams, fresult)
		} else {
			paramType = fmt.Sprintf("func(%s)", fparams)
		}
	case *ast.MapType:
		paramType = fmt.Sprintf("map[%s]%s", getType(t.Key, true), getType(t.Value, true))
	case *ast.ChanType:
		paramType = fmt.Sprintf("chan %s", getType(t.Value, true))
	}
	return
}

func getAccess(name string) (access string) {
	if idx := strings.LastIndex(name, "."); idx > -1 && idx < len(name) {
		name = name[idx+1:]
	}

	if ast.IsExported(name) {
		access = "public"
	} else {
		access = "private"
	}
	return
}

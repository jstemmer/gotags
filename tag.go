package main

import (
	"bytes"
	"strconv"
)

type Tag struct {
	Name    string
	File    string
	Address string
	Type    string
	Fields  map[string]string
}

func NewTag(name, file string, line int, tagtype string) *Tag {
	l := strconv.Itoa(line)
	return &Tag{
		Name:    name,
		File:    file,
		Address: l,
		Type:    tagtype,
		Fields:  map[string]string{"line": l},
	}
}

func (t Tag) String() string {
	var b bytes.Buffer

	b.WriteString(t.Name)
	b.WriteByte('\t')
	b.WriteString(t.File)
	b.WriteByte('\t')
	b.WriteString(t.Address)
	b.WriteString(";\"\t")
	b.WriteString(t.Type)

	for k, v := range t.Fields {
		b.WriteByte('\t')
		b.WriteString(k)
		b.WriteByte(':')
		b.WriteString(v)
	}

	return b.String()
}

package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Tag struct {
	Name    string
	File    string
	Address string
	Type    string
	Fields  map[string]string
}

func NewTag(name, file string, line int, tagtype string) Tag {
	l := strconv.Itoa(line)
	return Tag{
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
	b.WriteByte('\t')

	fields := make([]string, len(t.Fields))
	i := 0
	for k, v := range t.Fields {
		fields[i] = fmt.Sprintf("%s:%s", k, v)
		i++
	}

	sort.Sort(sort.StringSlice(fields))
	b.WriteString(strings.Join(fields, "\t"))

	return b.String()
}

package main

import (
	"testing"
)

func TestTag(t *testing.T) {
	tag := NewTag("tagname", "filename", 1, "type")
	tag.Fields["test"] = "value"

	expected := []string{
		"tagname\tfilename\t1;\"\ttype\tline:1\ttest:value",
		"tagname\tfilename\t1;\"\ttype\ttest:value\tline:1",
	}

	s := tag.String()
	if s != expected[0] && s != expected[1] {
		t.Errorf("Tag.String() == %s, want %s", s, expected[0])
	}
}

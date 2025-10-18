package main

import (
	"encoding/xml"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestHTMLSmoke(t *testing.T) {
	htmlString := mark2html(`# Hello
		world
	`)

	decoder := xml.NewDecoder(strings.NewReader(htmlString))

	for {
		_, err := decoder.Token()

		if errors.Is(err, io.EOF) {
			break
		}
		if err == nil {
			continue
		}

		t.Error(err)
		break
	}
}

func TestGlobber(t *testing.T) {
	mapfs := fstest.MapFS{
		"a/f":          &fstest.MapFile{},
		"b/f.md":       &fstest.MapFile{},
		"c/d/d/d/f.md": &fstest.MapFile{},
		"d.md":         &fstest.MapFile{},
		"e.md/d/d/f":   &fstest.MapFile{},
		"f.md/d/f.md":  &fstest.MapFile{},

		"g/1.md": &fstest.MapFile{},
		"g/2.md": &fstest.MapFile{},
		"g/3.md": &fstest.MapFile{},
	}

	assert := func(want string, arg string) {
		t.Helper()
		var gotSlice []string
		fs.WalkDir(mapfs, arg, globber{&gotSlice}.walkDirFunc)

		got := strings.Join(gotSlice, " ")
		if want != got {
			t.Errorf("want %s    got %s", want, got)
		}
	}

	assert("", "a")
	assert("b/f.md", "b")
	assert("c/d/d/d/f.md", "c")
	assert("d.md", "d.md")
	assert("", "e.md")
	assert("f.md/d/f.md", "f.md")

	assert("g/1.md g/2.md g/3.md", "g")
}

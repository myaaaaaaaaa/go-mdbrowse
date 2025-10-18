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

type errorFS struct {
	fs.FS
}

func (fsys errorFS) Open(name string) (fs.File, error) {
	if strings.Contains(name, "error") {
		return nil, errors.New("error file: " + name)
	}
	return fsys.FS.Open(name)
}

func TestGlobber(t *testing.T) {
	mapfs := errorFS{fstest.MapFS{
		"a/f":          &fstest.MapFile{},
		"b/f.md":       &fstest.MapFile{},
		"c/d/d/d/f.md": &fstest.MapFile{},
		"d.md":         &fstest.MapFile{},
		"e.md/d/d/f":   &fstest.MapFile{},
		"f.md/d/f.md":  &fstest.MapFile{},

		"g/1.md": &fstest.MapFile{},
		"g/2.md": &fstest.MapFile{},
		"g/3.md": &fstest.MapFile{},

		"h/1/1.md":       &fstest.MapFile{},
		"h/2/error/2.md": &fstest.MapFile{},
		"h/3/3.md":       &fstest.MapFile{},

		"i/error.md": &fstest.MapFile{},
	}}

	assert := func(want string, arg string) {
		t.Helper()
		var gotSlice []string
		err := fs.WalkDir(mapfs, arg, globber{&gotSlice}.walkDirFunc)
		if err != nil {
			t.Error(err)
		}

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

	assert("h/1/1.md h/3/3.md", "h")

	assert("i/error.md", "i")
}

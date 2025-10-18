package main

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
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

func TestFindMD(t *testing.T) {
	t.Chdir(t.TempDir())

	assert := func(want string, args ...string) {
		t.Helper()

		gotSlice := findMarkdownFiles(args)
		got := strings.Join(gotSlice, " ")
		if want != got {
			t.Errorf("want %s    got %s", want, got)
		}
	}

	os.CopyFS(".", fstest.MapFS{
		"a/f":          &fstest.MapFile{},
		"b/f.md":       &fstest.MapFile{},
		"c/d/d/d/f.md": &fstest.MapFile{},
		"d.md":         &fstest.MapFile{},
		"e.md/d/d/f":   &fstest.MapFile{},
		"f.md/d/f.md":  &fstest.MapFile{},

		"g/1.md": &fstest.MapFile{},
		"g/2.md": &fstest.MapFile{},
		"g/3.md": &fstest.MapFile{},
	})

	assert("")
	assert("", "a")
	assert("b/f.md", "b")
	assert("c/d/d/d/f.md", "c")
	assert("d.md", "d.md")
	assert("", "e.md")
	assert("f.md/d/f.md", "f.md")

	assert("g/1.md g/2.md g/3.md", "g")
}

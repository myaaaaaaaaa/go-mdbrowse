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

func assertEqual[T comparable](t *testing.T, want, got T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v    want %v", got, want)
	}
}

func TestHTMLSmoke(t *testing.T) {
	htmlString := mark2html(`# Hello
		world
	`)

	if !strings.Contains(htmlString, "<style") {
		t.Error("no <style> element found")
	}

	decoder := xml.NewDecoder(strings.NewReader(htmlString))

	var err error
	for err == nil {
		_, err = decoder.Token()
	}
	if !errors.Is(err, io.EOF) {
		t.Error(err)
	}
}

func TestTmpdir(t *testing.T) {
	assert := func(want, rundir, user string) {
		t.Helper()

		assertEqual(t, want, tmpdir(func(s string) string {
			switch s {
			case "XDG_RUNTIME_DIR":
				return rundir
			case "USER":
				return user
			}
			panic("invalid key: " + s)
		}))
	}

	assert("/tmp", "", "")
	assert("/tmp/usr", "", "usr")
	assert("/run/usr", "/run/usr", "")
	assert("/run/usr", "/run/usr", "usr")
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
		"g/error.md":   &fstest.MapFile{},

		"z/1.md": &fstest.MapFile{},
		"z/2.md": &fstest.MapFile{},
		"z/3.md": &fstest.MapFile{},

		"y/1/1.md":       &fstest.MapFile{},
		"y/2/error/2.md": &fstest.MapFile{},
		"y/3/3.md":       &fstest.MapFile{},
	}}

	assert := func(arg, want string) {
		t.Helper()
		var gotSlice []string
		err := fs.WalkDir(mapfs, arg, globber{&gotSlice}.walkDirFunc)
		if err != nil {
			t.Error(err)
		}

		assertEqual(t, want, strings.Join(gotSlice, " "))
	}

	assert("a", "")
	assert("b", "b/f.md")
	assert("c", "c/d/d/d/f.md")
	assert("d.md", "d.md")
	assert("e.md", "")
	assert("f.md", "f.md/d/f.md")
	assert("g", "g/error.md")

	assert("z", "z/1.md z/2.md z/3.md")

	assert("y", "y/1/1.md y/3/3.md")
}

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "embed"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func newMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
	)
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

//go:embed style.css
var css string

func mark2html(text string) string {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, `<meta name="viewport" content="width=device-width, initial-scale=1" />`)
	fmt.Fprintln(&buf, `<style>`)
	fmt.Fprintln(&buf, css)
	fmt.Fprintln(&buf, `</style>`)

	err := newMarkdown().Convert([]byte(text), &buf)

	if err != nil {
		rawText := "markdown error: " + err.Error() + "\n\n" + text
		return "<code>\n" + rawText + "\n</code>"
	}
	return buf.String()
}

type globber struct {
	files *[]string
}

func (g globber) walkDirFunc(p string, d os.DirEntry, err error) error {
	if err != nil {
		fmt.Println("skipping dir:", err)
		return fs.SkipDir
	}
	if !d.IsDir() && strings.HasSuffix(p, ".md") {
		*g.files = append(*g.files, p)
	}
	return nil
}

func findMarkdownFiles() (rt []string) {
	files := os.Args[1:]
	if len(files) == 0 {
		//files = []string{"."}
	}

	for _, file := range files {
		// For consistency
		file = filepath.Clean(file)

		err := filepath.WalkDir(file, globber{&rt}.walkDirFunc)
		if err != nil {
			panic(err)
		}
	}

	return
}

func printLink(f string) {
	f = must(filepath.Abs(f))
	f = strings.ReplaceAll(f, " ", "%20")

	fmt.Println("file://" + f)
}

func convertFile(markdownFile, htmlFile string) {
	data := must(os.ReadFile(markdownFile))
	data = []byte(mark2html(string(data)))
	must(0, os.WriteFile(htmlFile, data, 0600))

	fmt.Print(markdownFile, "    ")
	printLink(htmlFile)
}

func mktemp() string {
	basedir := os.Getenv("XDG_RUNTIME_DIR")
	if basedir == "" {
		basedir = "/tmp"

		user := os.Getenv("USER")
		if user != "" {
			basedir += "/" + user
		}
	}

	basedir += "/markdown/"
	basedir += fmt.Sprint(time.Now().Unix())
	return basedir
}

//go:embed go.mod
var gomod string

func main() {
	{
		pkgName := gomod
		pkgName, _, _ = strings.Cut(pkgName, "\n")
		_, pkgName, _ = strings.Cut(pkgName, " ")
		pkgName = strings.TrimSpace(pkgName)

		fmt.Println()
		fmt.Printf("\tgo install %s@latest\n", pkgName)
		fmt.Println()
	}

	markdownFiles := findMarkdownFiles()

	outDir := mktemp()
	must(0, os.MkdirAll(outDir, 0700))

	i := 0
	for _, markdownFile := range markdownFiles {
		i++
		htmlFile := fmt.Sprintf("%s/%05d.htm", outDir, i)

		convertFile(markdownFile, htmlFile)
	}

	var tmpFiles []string
	{
		tmpFiles = append(tmpFiles, must(filepath.Glob("/tmp/*.md"))...)
		rundir := os.Getenv("XDG_RUNTIME_DIR")
		if rundir != "" {
			tmpFiles = append(tmpFiles, must(filepath.Glob(rundir+"/*.md"))...)
		}
	}
	for _, tmpFile := range tmpFiles {
		htmlFile := tmpFile + ".htm"

		_, statErr := os.Stat(htmlFile)
		if statErr == nil {
			continue
		}

		convertFile(tmpFile, htmlFile)
	}
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "embed"

	"github.com/yuin/goldmark"
)

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

//go:embed markdown.css
var css string

func mark2html(text string) string {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, `<meta name="viewport" content="width=device-width, initial-scale=1" />`)
	fmt.Fprintln(&buf, `<style>`)
	fmt.Fprintln(&buf, css)
	fmt.Fprintln(&buf, `</style>`)

	err := goldmark.Convert([]byte(text), &buf)

	if err != nil {
		rawText := "markdown error: " + err.Error() + "\n\n" + text
		return "<code>\n" + rawText + "\n</code>"
	}
	return buf.String()
}

func findMarkdownFiles() (rt []string) {
	files := os.Args[1:]
	if len(files) == 0 {
		files = []string{"."}
	}

	for _, file := range files {
		// For consistency
		file = filepath.Clean(file)

		err := filepath.WalkDir(file, func(p string, d os.DirEntry, err error) error {
			if !d.IsDir() && strings.HasSuffix(p, ".md") {
				rt = append(rt, p)
			}
			return err
		})

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
		url := gomod
		url, _, _ = strings.Cut(url, "\n")
		_, url, _ = strings.Cut(url, " ")
		url = strings.TrimSpace(url)

		fmt.Println()
		fmt.Printf("\tgo install %s@latest\n", url)
		fmt.Println()
	}

	markdownFiles := findMarkdownFiles()

	outDir := mktemp()
	must(0, os.MkdirAll(outDir, 0700))

	i := 0
	for _, markdownFile := range markdownFiles {
		i++
		htmlFile := fmt.Sprintf("%s/%05d.htm", outDir, i)

		data := must(os.ReadFile(markdownFile))
		data = []byte(mark2html(string(data)))
		must(0, os.WriteFile(htmlFile, data, 0600))

		fmt.Print(markdownFile, "    ")
		printLink(htmlFile)
	}
}

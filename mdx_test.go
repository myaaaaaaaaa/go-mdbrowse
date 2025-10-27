package main

import (
	"regexp"
	"strings"
	"testing"

	_ "embed"
)

func TestTokenize(t *testing.T) {
	indentRe := regexp.MustCompile(`(?m)^\t+`)
	trim := func(s string) string {
		s = indentRe.ReplaceAllString(s, "")
		s = strings.Trim(s, "\n")
		return s
	}
	checkTokenize := func(want string) []string {
		t.Helper()

		want = trim(want)

		tokens := tokenizeHeadings(want)
		for _, token := range tokens {
			assertEqual(t, false, token == "")
		}
		got := strings.Join(tokens, "")
		assertEqual(t, want, got)

		return tokens
	}

	assertCut := func(md, want string) {
		t.Helper()

		want = trim(want)
		md = trim(md)

		got := checkTokenize(md)
		assertEqual(t, want, strings.Join(got, "|"))
	}

	assertCut(`
		# Hello world!
		Who am I to say...?

		## Goodbye sky!
	`, `
		# Hello world!
		|Who am I to say...?
		|
		|## Goodbye sky!
	`)

	assertCut(`
		# Hello world!
		~~~# h1
		~
		~~
		## h2
		~~~
		# h1
		## h2
		text
		~~~
		...
	`, `
		# Hello world!
		|~~|~# h1
		~
		~~
		## h2
		~~|~
		|# h1
		|## h2
		|text
		|~~|~
		...
	`)

	assertCut(`
		~~~~
		~~~
		~~~~
		aa
		~~~~~
		~~~~
		~~
		~~~~~
		aa
		~~~~~~
	`, `
		~~~|~
		~~~
		~~~|~
		|aa
		|~~~~|~
		~~~~
		~~
		~~~~|~
		|aa
		|~~~~~~
	`)
}

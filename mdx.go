package main

import (
	"strings"
)

func consumeToken(rt *[]string, md *string, f func(byte) (int, bool)) string {
	i := 0
	for ; i < len(*md); i++ {
		n, ok := f((*md)[i])
		i += n

		if !ok {
			break
		}
	}

	head := ""
	head, *md = (*md)[:i], (*md)[i:]
	if head != "" {
		*rt = append(*rt, head)
	}
	return head
}

func tokenizeHeadings(md string) []string {
	var rt []string

	for md != "" {
		if strings.HasPrefix(md, "```") || strings.HasPrefix(md, "~~~") {
			quote := md[0]
			quoteStr := consumeToken(&rt, &md, func(c byte) (int, bool) {
				if c != quote {
					return -1, false
				}
				return 0, true
			})

			count := 0

			consumeToken(&rt, &md, func(c byte) (int, bool) {
				if c == quote {
					count++
				} else {
					count = 0
				}

				if count > len(quoteStr) {
					return 0, false
				}
				return 0, true
			})

			continue
		}

		consumeToken(&rt, &md, func(c byte) (int, bool) {
			if c == '\n' {
				return 1, false
			}
			return 0, true
		})
	}

	return rt
}

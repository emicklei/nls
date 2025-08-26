package main

import "strings"

type Entry struct {
	Language    string
	Key         string
	Text        string
	Description string
}

func (e Entry) Replacements() int {
	return strings.Count(e.Text, "{{.")
}

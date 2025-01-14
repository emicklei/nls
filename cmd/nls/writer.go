package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func writeEntries(entries []Entry, dir string) error {
	mme := entriesToMap(entries)
	for lang, langMap := range mme {
		if err := writeLangFile(lang, langMap, dir); err != nil {
			return err
		}
	}
	return nil
}

func entriesToMap(entries []Entry) map[string]map[string]Entry {
	msg := map[string]map[string]Entry{}
	for _, each := range entries {
		var langMap map[string]Entry
		if lm, ok := msg[each.Language]; !ok {
			langMap = map[string]Entry{}
			msg[each.Language] = langMap
		} else {
			langMap = lm
		}
		entry, ok := langMap[each.Key]
		if !ok {
			langMap[each.Key] = each
		} else {
			// only overwrite if the text is not empty
			if each.Text != "" {
				entry.Text = each.Text
				langMap[each.Key] = entry
			}
		}
	}
	return msg
}

func writeLangFile(lang string, langMap map[string]Entry, dir string) error {
	fileName := filepath.Join(dir, lang, "messages.yaml")
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	// collect entries
	entries := []Entry{}
	for _, each := range langMap {
		entries = append(entries, each)
	}
	// sort by key
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	for _, each := range entries {
		fmt.Fprintf(out, "%s: ", each.Key)
		writeYAMLString(out, each.Text)
	}
	return nil
}

var quoteit = "{}[],&*#?|-<>=!%@"

func writeYAMLString(w io.Writer, s string) {
	if s == "" {
		fmt.Fprintln(w)
		return
	}
	if strings.Contains(s, "\n") {
		fmt.Fprintln(w, "|")
		lines := strings.Split(s, "\n")
		for i, line := range lines {
			if line == "" && i == len(lines)-1 {
				continue
			}
			fmt.Fprintf(w, "  %s\n", line)
		}
		return
	}
	if strings.ContainsAny(s, quoteit) {
		fmt.Fprintf(w, "'%s'\n", s)
		return
	}
	fmt.Fprintf(w, "%s\n", s)
}

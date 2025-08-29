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
				// keep comment of existing
				each.Comment = entry.Comment
				langMap[each.Key] = each
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
		if each.Comment != "" {
			fmt.Fprintln(out, each.Comment)
		}
		fmt.Fprintf(out, "%s: ", each.Key)
		if each.Description != "" {
			writeNestedYAMLString(out, each.Text, each.Description)
		} else {
			writeYAMLString(out, each.Text)
		}
	}
	return nil
}

var quoteit = "{}[],&*#?|-<>=!%@"

// write the key using a nested format:
// key:
//
//	msg: value
//	desc: explanation of the context in which the value is used
func writeNestedYAMLString(w io.Writer, msg string, desc string) {
	fmt.Fprintln(w)

	fmt.Fprint(w, "  msg: ")
	if msg == "" {
		fmt.Fprintln(w)
	} else if strings.Contains(msg, "\n") {
		fmt.Fprintln(w, "|")
		lines := strings.Split(msg, "\n")
		for i, line := range lines {
			if line == "" && i == len(lines)-1 {
				continue
			}
			fmt.Fprintf(w, "    %s\n", line)
		}
	} else if strings.ContainsAny(msg, quoteit) {
		fmt.Fprintf(w, "'%s'\n", msg)
	} else {
		fmt.Fprintf(w, "%s\n", msg)
	}

	fmt.Fprint(w, "  desc: ")
	if desc == "" {
		fmt.Fprintln(w)
	} else if strings.Contains(desc, "\n") {
		fmt.Fprintln(w, "|")
		lines := strings.Split(desc, "\n")
		for i, line := range lines {
			if line == "" && i == len(lines)-1 {
				continue
			}
			fmt.Fprintf(w, "    %s\n", line)
		}
	} else if strings.ContainsAny(desc, quoteit) {
		fmt.Fprintf(w, "'%s'\n", desc)
	} else {
		fmt.Fprintf(w, "%s\n", desc)
	}
}

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

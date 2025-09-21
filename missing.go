package nls

import (
	"fmt"
	"strings"
)

var Missing map[string]Fallback = map[string]Fallback{}

// for storing missing ones
type Fallback struct {
	Lang string
	Key  string
	Msg  string
}

func addMissing(lang, key, msg string) {
	Missing[fmt.Sprintf("%s::%s", lang, key)] = Fallback{Lang: lang, Key: key, Msg: msg}
}

func ReportMissing() string {
	report := new(strings.Builder)
	// build by lang
	byLang := map[string][]Fallback{}
	for _, e := range Missing {
		byLang[e.Lang] = append(byLang[e.Lang], e)
	}
	for lang, entries := range byLang {
		fmt.Fprintf(report, "%s:\n", lang)
		for _, entry := range entries {
			fmt.Fprintf(report, "\t%s:\n\t\tmsg: %s\n\t\tdesc:\n", entry.Key, entry.Msg)
		}
	}
	return report.String()
}

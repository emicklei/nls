package nls

import (
	"fmt"
	"strings"
	"text/template"
)

var Missing map[string]Fallback = map[string]Fallback{}
var trackMissing = true

// for storing missing ones
type Fallback struct {
	Lang string
	Key  string
	Msg  string
}

func addMissing(lang, key, msg string) {
	if !trackMissing {
		return
	}
	Missing[lang+"::"+key] = Fallback{Lang: lang, Key: key, Msg: msg}
}

// EnableMissingTracking controls whether missing keys are recorded.
func EnableMissingTracking(enabled bool) {
	trackMissing = enabled
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
			fmt.Fprintf(report, "\t%s:\n", entry.Key)
			fmt.Fprintf(report, "\t\tmsg: %s\n", entry.Msg)
			fmt.Fprintf(report, "\t\tdesc:\n")
		}
	}
	return report.String()
}

// Register is called from generated code.
func Register(catalog map[string]*template.Template, key string, templateSource string) {
	catalog[key] = template.Must(template.New(key).Parse(templateSource))
}

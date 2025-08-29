package nls

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type Localizer interface {
	// Get returns the text associated with a key for using the available languages
	// It returns an empty string if none of the languages have a (non-empty) value for the key and no fallback is provided.
	Get(key string, fallback ...string) string
	// Format returns the text after applying substitutions using the key(string) and value pairs.
	// Returns an empty string if there no such key.
	Format(key string, kv ...any) string
	// Replaced returns the text after applying substitutions using the replacements.
	// Returns an empty string if there no such key.
	Replaced(key string, replacements ...map[string]any) string
	// ReportMissing returns a report of all missing translations.
	ReportMissing() string
}

type localizer struct {
	catalog   map[string]*template.Template
	languages []string
	missing   map[string]string
}

func NewLocalizer(catalog map[string]*template.Template, languages ...string) Localizer {
	return localizer{catalog: catalog, languages: languages, missing: map[string]string{}}
}

func (l localizer) findTemplate(key string) *template.Template {
	for _, lang := range l.languages {
		mapkey := fmt.Sprintf("%s.%s", lang, key)
		if tmpl, ok := l.catalog[mapkey]; ok {
			return tmpl
		}
	}
	return nil
}

func (l localizer) ReportMissing() string {
	report := new(strings.Builder)
	for k, v := range l.missing {
		fmt.Fprintf(report, "%s:\n\tmsg: %s\n\tdesc:\n", k, v)
	}
	return report.String()
}

// Get returns the text associated with a key for using the available languages
// It returns an empty string if none of the languages have a (non-empty) value for the key and no fallback is provided.
func (l localizer) Get(key string, fallback ...string) string {
	tmpl := l.findTemplate(key)
	if tmpl == nil {
		if len(fallback) > 0 {
			l.missing[key] = fallback[0]
			return fallback[0]
		}
		return ""
	}
	buf := new(bytes.Buffer)
	// execute with no data
	if err := tmpl.Execute(buf, nil); err != nil {
		return err.Error()
	}
	msg := buf.String()
	if msg == "" {
		if len(fallback) > 0 {
			l.missing[key] = fallback[0]
			return fallback[0]
		}
	}
	return msg
}

// Format returns the text after applying substitutions using the key(string) and value pairs.
// Returns an empty string if there no such key.
func (l localizer) Format(key string, kv ...any) string {
	params := map[string]any{}
	for i := 0; i < len(kv); i += 2 {
		k := kv[i]
		if ks, ok := k.(string); ok {
			params[ks] = kv[i+1]
		} else {
			return "bad arguments: Format expects [string,any] pairs"
		}
	}
	return l.Replaced(key, params)
}

// Replaced returns the text after applying substitutions using the replacements.
// Returns an empty string if there no such key.
func (l localizer) Replaced(key string, replacements ...map[string]any) string {
	tmpl := l.findTemplate(key)
	if tmpl == nil {
		return ""
	}
	var data any
	if len(replacements) > 0 {
		data = replacements[0]
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return err.Error()
	}
	return buf.String()
}

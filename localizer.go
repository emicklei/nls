package nls

import (
	"fmt"
	"strings"
	"text/template"
)

type Localizer struct {
	catalog   map[string]string
	languages []string
}

func NewLocalizer(catalog map[string]string, languages ...string) Localizer {
	return Localizer{catalog: catalog, languages: languages}
}

// Get returns the text associated with a key for using the available languages
// It returns an empty string if none of the languages have a (non-empty) value for the key and no fallback is provided.
func (l Localizer) Get(key string, fallback ...string) string {
	for _, lang := range l.languages {
		mapkey := fmt.Sprintf("%s.%s", lang, key)
		if v, ok := l.catalog[mapkey]; ok && len(v) > 0 {
			return v
		}
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

// Format returns the text after applying substitutions using the key(string) and value pairs.
// Returns an empty string if there no such key.
func (l Localizer) Format(key string, kv ...any) string {
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
func (l Localizer) Replaced(key string, replacements ...map[string]any) string {
	tmpl := l.Get(key)
	// If no replacements are provided, return the template as is.
	if len(replacements) == 0 {
		return tmpl
	}
	// If the tmpl doesn't have any substitutions, no need to template.Execute.
	// Note: this optimization is removed because it prevents detection of malformed templates.
	// if !strings.Contains(tmpl, "}}") {
	// 	return tmpl
	// }
	replacer, err := template.New("replacer").Parse(tmpl)
	if err != nil {
		return err.Error()
	}
	buf := new(strings.Builder)
	if err := replacer.Execute(buf, replacements[0]); err != nil {
		return err.Error()
	}
	return buf.String()
}

package nls

import (
	"fmt"
	"strings"
	"text/template"
)

const (
	M_hello = "hello"
	M_sky   = "sky"
	M_world = "world"
	M_sea   = "sea"
	M_cats  = "cats"
)

var messages = map[string]string{
	"en." + M_hello: "Hello",
	"nl." + M_hello: "Hallo",
	"en." + M_world: "World",
	"nl." + M_world: "Wereld",
	"nl." + M_sea:   "{{.name}} zee",
	"nl." + M_cats:  "{{.count}} {{- if gt .count 1}} katten{{- else}} kat{{- end}}",
}

type Localizer struct {
	languages []string
}

func New(languages ...string) Localizer {
	return Localizer{languages: languages}
}

func (l Localizer) Get(key string, fallback ...string) string {
	for _, lang := range l.languages {
		mapkey := fmt.Sprintf("%s.%s", lang, key)
		if v, ok := messages[mapkey]; ok {
			return v
		}
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}
func (l Localizer) Replaced(key string, replacements ...map[string]any) string {
	tmpl := l.Get(key)
	if len(replacements) == 0 {
		return tmpl
	}
	// If no replacements are provided, return the template as is.
	if len(replacements) == 0 {
		return tmpl
	}
	// If the tmpl doesn't have any substitutions, no need to template.Execute.
	if !strings.Contains(tmpl, "}}") {
		return tmpl
	}
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

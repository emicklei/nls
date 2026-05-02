package nls

import (
	"bytes"
	"strings"
	"text/template"
	"text/template/parse"

	"golang.org/x/text/language"
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
}

type localizer struct {
	catalog   map[string]*template.Template
	languages []string // at least one language is present
	resolved  map[string]*template.Template
	plain     map[string]string
}

func NewLocalizer(catalog map[string]*template.Template, languages ...string) Localizer {
	if len(languages) == 0 {
		languages = append(languages, language.English.String())
	}
	resolved, plain := resolveCatalog(catalog, languages)
	return localizer{catalog: catalog, languages: languages, resolved: resolved, plain: plain}
}

func (l localizer) findTemplate(key string) *template.Template {
	return l.resolved[key]
}

func resolveCatalog(catalog map[string]*template.Template, languages []string) (map[string]*template.Template, map[string]string) {
	resolved := map[string]*template.Template{}
	plain := map[string]string{}
	if len(catalog) == 0 || len(languages) == 0 {
		return resolved, plain
	}
	priorities := make(map[string]int, len(languages))
	for i, lang := range languages {
		if _, exists := priorities[lang]; !exists {
			priorities[lang] = i
		}
	}
	selected := map[string]int{}
	for fullKey, tmpl := range catalog {
		dot := strings.IndexByte(fullKey, '.')
		if dot <= 0 || dot == len(fullKey)-1 {
			continue
		}
		lang := fullKey[:dot]
		priority, ok := priorities[lang]
		if !ok {
			continue
		}
		key := fullKey[dot+1:]
		if current, exists := selected[key]; exists && priority >= current {
			continue
		}
		selected[key] = priority
		resolved[key] = tmpl
		if text, ok := staticTemplateText(tmpl); ok {
			plain[key] = text
		} else {
			delete(plain, key)
		}
	}
	return resolved, plain
}

func staticTemplateText(tmpl *template.Template) (string, bool) {
	if tmpl == nil || tmpl.Tree == nil || tmpl.Tree.Root == nil {
		return "", false
	}
	nodes := tmpl.Tree.Root.Nodes
	if len(nodes) == 0 {
		return "", true
	}
	if len(nodes) == 1 {
		if textNode, ok := nodes[0].(*parse.TextNode); ok {
			return string(textNode.Text), true
		}
	}
	return "", false
}

// Get returns the text associated with a key for using the available languages
// It returns an empty string if none of the languages have a (non-empty) value for the key and no fallback is provided.
func (l localizer) Get(key string, fallback ...string) string {
	if msg, ok := l.plain[key]; ok {
		if msg == "" {
			if len(fallback) > 0 {
				addMissing(l.languages[0], key, fallback[0])
				return fallback[0]
			}
			addMissing(l.languages[0], key, "")
		}
		return msg
	}
	tmpl := l.findTemplate(key)
	if tmpl == nil {
		if len(fallback) > 0 {
			addMissing(l.languages[0], key, fallback[0])
			return fallback[0]
		}
		addMissing(l.languages[0], key, "")
		return key
	}
	buf := new(bytes.Buffer)
	// execute with no data
	_ = tmpl.Execute(buf, nil)
	msg := buf.String()
	if msg == "" {
		if len(fallback) > 0 {
			addMissing(l.languages[0], key, fallback[0])
			return fallback[0]
		}
		addMissing(l.languages[0], key, "")
	}
	return msg
}

// Format returns the text after applying substitutions using the key(string) and value pairs.
// Returns an empty string if there no such key.
func (l localizer) Format(key string, kv ...any) string {
	if len(kv)%2 != 0 {
		return "bad arguments: Format expects [string,any] pairs"
	}
	params := make(map[string]any, len(kv)/2)
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
	if msg, ok := l.plain[key]; ok {
		return msg
	}
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

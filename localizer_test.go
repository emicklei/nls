package nls

import (
	"testing"
)

func TestGet(t *testing.T) {
	cat := map[string]string{
		"en.hello": "world",
		"nl.hello": "wereld",
		"en.empty": "",
		"nl.empty": "no value",
	}
	l := NewLocalizer(cat, "nl", "en")
	if got, want := l.Get("hello"), "wereld"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	l = NewLocalizer(cat, "en", "nl")
	if got, want := l.Get("hello"), "world"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	l = NewLocalizer(cat, "fr", "en")
	if got, want := l.Get("hello"), "world"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	l = NewLocalizer(cat, "fr")
	if got, want := l.Get("hello"), ""; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	l = NewLocalizer(cat, "fr")
	if got, want := l.Get("hello", "fallback"), "fallback"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	l = NewLocalizer(cat, "en")
	if got, want := l.Get("unknown"), ""; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	l = NewLocalizer(cat, "en", "nl")
	if got, want := l.Get("empty"), "no value"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestFormat(t *testing.T) {
	cat := map[string]string{
		"en.template": "this is a {{.what}}",
		"en.multi":    "this is a {{.what}} and {{.who}}",
		"en.no_subst": "this is a test",
	}
	l := NewLocalizer(cat, "en")
	if got, want := l.Format("template", "what", "test"), "this is a test"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Format("multi", "what", "test", "who", "me"), "this is a test and me"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Format("unknown"), ""; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Format("no_subst", "what", "test"), "this is a test"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Format("template", 1, "test"), "bad arguments: Format expects [string,any] pairs"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestReplaced(t *testing.T) {
	cat := map[string]string{
		"en.template":       "this is a {{.what}}",
		"en.no_subst":       "this is a test",
		"en.invalid_tmpl":   "this is a {{.what",
		"en.no_repl_needed": "no replacements",
	}
	l := NewLocalizer(cat, "en")
	if got, want := l.Replaced("template", map[string]any{"what": "test"}), "this is a test"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Replaced("unknown"), ""; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Replaced("no_subst", map[string]any{"what": "test"}), "this is a test"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Replaced("invalid_tmpl", map[string]any{"what": "test"}), "template: replacer:1: unclosed action"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Replaced("no_repl_needed"), "no replacements"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

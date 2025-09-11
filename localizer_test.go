package nls

import (
	"strings"
	"testing"
	"text/template"
)

func mustTemplate(s string) *template.Template {
	return template.Must(template.New("").Parse(s))
}

func TestNewLocalizer(t *testing.T) {
	l := NewLocalizer(nil)
	if got, want := l.(localizer).languages[0], "en"; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
}

func TestGet(t *testing.T) {
	cat := map[string]*template.Template{
		"en.hello":      mustTemplate("world"),
		"nl.hello":      mustTemplate("wereld"),
		"en.empty":      mustTemplate(""),
		"nl.empty":      mustTemplate("no value"),
		"en.exec_error": mustTemplate("{{.error}}"),
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
	if got, want := l.Get("empty"), ""; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := l.Get("empty", "fallback"), "fallback"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	// test missing report for empty
	l.Get("empty")
	if got, want := l.ReportMissing(), "en:\n\tempty:\n\t\tmsg: \n\t\tdesc:\n"; !strings.Contains(got, want) {
		t.Errorf("got [%s] should contain [%s]", got, want)
	}
	l = NewLocalizer(cat, "en")
	if got := l.Get("exec_error"); got == "" {
		t.Error("expected an error string")
	}
	// test missing report for empty
	l.Get("empty")
}

func TestFormat(t *testing.T) {
	cat := map[string]*template.Template{
		"en.template": mustTemplate("this is a {{.what}}"),
		"en.multi":    mustTemplate("this is a {{.what}} and {{.who}}"),
		"en.no_subst": mustTemplate("this is a test"),
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
	cat := map[string]*template.Template{
		"en.template":       mustTemplate("this is a {{.what}}"),
		"en.no_subst":       mustTemplate("this is a test"),
		"en.no_repl_needed": mustTemplate("no replacements"),
		"en.exec_error":     mustTemplate("{{index .A 1}}"),
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
	if got, want := l.Replaced("no_repl_needed"), "no replacements"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	// trigger an error during template execution
	if got, want := l.Replaced("exec_error", map[string]any{"A": []string{}}), `error calling index: index out of range`; !strings.Contains(got, want) {
		t.Errorf("got [%v] want to contain [%v]", got, want)
	}
}

func TestMissing(t *testing.T) {
	cat := map[string]*template.Template{
		"en.empty": mustTemplate(""),
	}
	l := NewLocalizer(cat, "en")
	l.Get("absent", "value")
	l.Get("empty")
	report := l.ReportMissing()
	if !strings.Contains(report, "absent") {
		t.Error("missing absent entry")
	}
	if !strings.Contains(report, "empty") {
		t.Error("missing empty entry")
	}
}

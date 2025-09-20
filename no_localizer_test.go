package nls

import "testing"

func TestNoLocalizer(t *testing.T) {
	nl := NoLocalizer{}
	if got, want := nl.Get("hello"), "hello"; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
	if got, want := nl.Get("hello", "world"), "world"; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
	if got, want := nl.Format("hello %s", "world"), "hello %s"; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
	if got, want := nl.Replaced("hello %s", map[string]any{"arg": "world"}), "hello %s"; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
	if got, want := nl.ReportMissing(), ""; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
}
func TestDefaultLocalizer(t *testing.T) {
	if _, ok := DefaultLocalizer().(Localizer); !ok {
		t.Fatal("must be Localizer")
	}
}

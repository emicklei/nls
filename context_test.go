package nls

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	if _, ok := LocalizerFromContext(ctx).(localizer); ok {
		t.Fatal("must not have localizer")
	}
	SetDefault(NewLocalizer(nil, "en"))
	if _, ok := LocalizerFromContext(ctx).(localizer); !ok {
		t.Fatal("must have localizer")
	}
	ctx = context.WithValue(ctx, localizerKey, NewLocalizer(nil, "nl"))
	if l, ok := LocalizerFromContext(ctx).(localizer); !ok {
		t.Fatal("must have localizer")
	} else {
		if l.languages[0] != "nl" {
			t.Fatal("wrong localizer")
		}
	}
	ctx = ContextWithLocalizer(context.Background(), NewLocalizer(nil, "fr"))
	if l, ok := LocalizerFromContext(ctx).(localizer); !ok {
		t.Fatal("must have localizer")
	} else {
		if l.languages[0] != "fr" {
			t.Fatal("wrong localizer")
		}
	}
}

func TestTranslations(t *testing.T) {
	ctx := ContextWithLocalizer(context.Background(), NewLocalizer(nil, "fr"))
	if got, want := Get(ctx, "hello"), "hello"; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
	if got, want := Replaced(ctx, "hello %s", map[string]any{"arg": "world"}), ""; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
	if got, want := Format(ctx, "hello %d", "arg", 1), ""; got != want {
		t.Errorf("got [%s] want [%s]", got, want)
	}
}

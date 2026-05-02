package nls

import (
	"context"
	"testing"
	"text/template"
)

func benchmarkCatalog() map[string]*template.Template {
	return map[string]*template.Template{
		"en.hello":   mustTemplate("hello world"),
		"nl.hello":   mustTemplate("hallo wereld"),
		"en.multi":   mustTemplate("this is a {{.what}} and {{.who}}"),
		"en.empty":   mustTemplate(""),
		"en.dynamic": mustTemplate("hi {{.name}}"),
	}
}

func BenchmarkGetHit(b *testing.B) {
	EnableMissingTracking(false)
	b.Cleanup(func() { EnableMissingTracking(true) })
	l := NewLocalizer(benchmarkCatalog(), "en", "nl")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Get("hello")
	}
}

func BenchmarkGetMissWithFallback(b *testing.B) {
	EnableMissingTracking(false)
	b.Cleanup(func() { EnableMissingTracking(true) })
	l := NewLocalizer(benchmarkCatalog(), "en", "nl")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Get("unknown", "fallback")
	}
}

func BenchmarkReplacedDynamic(b *testing.B) {
	EnableMissingTracking(false)
	b.Cleanup(func() { EnableMissingTracking(true) })
	l := NewLocalizer(benchmarkCatalog(), "en", "nl")
	replacements := map[string]any{"name": "copilot"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Replaced("dynamic", replacements)
	}
}

func BenchmarkFormatSmallKV(b *testing.B) {
	EnableMissingTracking(false)
	b.Cleanup(func() { EnableMissingTracking(true) })
	l := NewLocalizer(benchmarkCatalog(), "en", "nl")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Format("multi", "what", "test", "who", "me")
	}
}

func BenchmarkContextGetWrapper(b *testing.B) {
	EnableMissingTracking(false)
	b.Cleanup(func() { EnableMissingTracking(true) })
	ctx := ContextWithLocalizer(context.Background(), NewLocalizer(benchmarkCatalog(), "en", "nl"))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Get(ctx, "hello")
	}
}

func BenchmarkContextGetLocalizerOnce(b *testing.B) {
	EnableMissingTracking(false)
	b.Cleanup(func() { EnableMissingTracking(true) })
	ctx := ContextWithLocalizer(context.Background(), NewLocalizer(benchmarkCatalog(), "en", "nl"))
	l := LocalizerFromContext(ctx)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Get("hello")
	}
}

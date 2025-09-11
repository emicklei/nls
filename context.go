package nls

import (
	"context"
)

var defaultLocalizer Localizer = NoLocalizer{}

// SetDefault set the default localizer which is used when none is present in a context.Context.
func SetDefault(localizer Localizer) {
	defaultLocalizer = localizer
}

// DefaultLocalizer returns the default one.
func DefaultLocalizer() Localizer { return defaultLocalizer }

var localizerKey = struct{ _ string }{}

// ContextWithLocalizer returns a new context that holds a Localizer.
func ContextWithLocalizer(ctx context.Context, localizer Localizer) context.Context {
	return context.WithValue(ctx, localizerKey, localizer)
}

// LocalizerFromContext returns the localizer from the context or the default is absent.
func LocalizerFromContext(ctx context.Context) Localizer {
	l, ok := ctx.Value(localizerKey).(Localizer)
	if ok {
		return l
	}
	return defaultLocalizer
}

// Get a localized string by its message ID, with an optional fallback.
func Get(ctx context.Context, messageID string, fallback ...string) string {
	return LocalizerFromContext(ctx).Get(messageID, fallback...)
}

// Replaced returns a localized string by its message ID, with optional replacements.
func Replaced(ctx context.Context, messageID string, replacements ...map[string]any) string {
	return LocalizerFromContext(ctx).Replaced(messageID, replacements...)
}

// Format returns a localized string by its message ID, with optional key-value pairs.
func Format(ctx context.Context, messageID string, kv ...any) string {
	return LocalizerFromContext(ctx).Format(messageID, kv...)
}

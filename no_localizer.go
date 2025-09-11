package nls

var _ Localizer = NoLocalizer{}

type NoLocalizer struct{}

// Get returns the first fallback if it is provided, otherwise it returns an empty string.
func (n NoLocalizer) Get(key string, fallback ...string) string {
	if len(fallback) > 0 {
		return fallback[0]
	}
	return key
}

// Format returns an empty string.
func (n NoLocalizer) Format(key string, kv ...any) string {
	return key
}

// Replaced returns an empty string.
func (n NoLocalizer) Replaced(key string, replacements ...map[string]any) string {
	return key
}

// ReportMissing returns an empty string.
func (n NoLocalizer) ReportMissing() string {
	return ""
}

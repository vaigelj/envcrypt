package envfile

import (
	"strings"
	"unicode"
)

// SanitizeOption configures the Sanitize function.
type SanitizeOption func(*sanitizeConfig)

type sanitizeConfig struct {
	stripControlChars bool
	normalizeNewlines bool
	trimQuotes        bool
	removeNullBytes   bool
}

// WithStripControlChars removes non-printable control characters from values.
func WithStripControlChars() SanitizeOption {
	return func(c *sanitizeConfig) { c.stripControlChars = true }
}

// WithNormalizeNewlines replaces \r\n and bare \r with \n in values.
func WithNormalizeNewlines() SanitizeOption {
	return func(c *sanitizeConfig) { c.normalizeNewlines = true }
}

// WithTrimQuotes removes surrounding single or double quotes from values.
func WithTrimQuotes() SanitizeOption {
	return func(c *sanitizeConfig) { c.trimQuotes = true }
}

// WithRemoveNullBytes strips null bytes (\x00) from values.
func WithRemoveNullBytes() SanitizeOption {
	return func(c *sanitizeConfig) { c.removeNullBytes = true }
}

// Sanitize applies the given options to clean up entry values in place,
// returning a new slice of sanitized entries.
func Sanitize(entries []Entry, opts ...SanitizeOption) []Entry {
	cfg := &sanitizeConfig{}
	for _, o := range opts {
		o(cfg)
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		e.Value = sanitizeValue(e.Value, cfg)
		result[i] = e
	}
	return result
}

func sanitizeValue(v string, cfg *sanitizeConfig) string {
	if cfg.removeNullBytes {
		v = strings.ReplaceAll(v, "\x00", "")
	}
	if cfg.normalizeNewlines {
		v = strings.ReplaceAll(v, "\r\n", "\n")
		v = strings.ReplaceAll(v, "\r", "\n")
	}
	if cfg.stripControlChars {
		v = strings.Map(func(r rune) rune {
			if unicode.IsControl(r) && r != '\n' && r != '\t' {
				return -1
			}
			return r
		}, v)
	}
	if cfg.trimQuotes {
		if len(v) >= 2 {
			if (v[0] == '"' && v[len(v)-1] == '"') ||
				(v[0] == '\'' && v[len(v)-1] == '\'') {
				v = v[1 : len(v)-1]
			}
		}
	}
	return v
}

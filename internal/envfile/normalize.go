package envfile

import (
	"strings"
	"unicode"
)

// NormalizeOption configures how normalization is applied.
type NormalizeOption func(*normalizeConfig)

type normalizeConfig struct {
	upperKeys   bool
	trimValues  bool
	quoteValues bool
	removeEmpty bool
}

// WithUpperKeys converts all keys to uppercase.
func WithUpperKeys() NormalizeOption {
	return func(c *normalizeConfig) { c.upperKeys = true }
}

// WithTrimValues trims leading/trailing whitespace from values.
func WithTrimValues() NormalizeOption {
	return func(c *normalizeConfig) { c.trimValues = true }
}

// WithQuoteValues wraps values that contain spaces in double quotes.
func WithQuoteValues() NormalizeOption {
	return func(c *normalizeConfig) { c.quoteValues = true }
}

// WithRemoveEmpty drops entries whose value is empty after trimming.
func WithRemoveEmpty() NormalizeOption {
	return func(c *normalizeConfig) { c.removeEmpty = true }
}

// Normalize applies the requested normalization passes to entries.
func Normalize(entries []Entry, opts ...NormalizeOption) []Entry {
	cfg := &normalizeConfig{}
	for _, o := range opts {
		o(cfg)
	}

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if cfg.upperKeys {
			e.Key = strings.ToUpper(e.Key)
		}
		if cfg.trimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if cfg.removeEmpty && e.Value == "" {
			continue
		}
		if cfg.quoteValues && needsQuoting(e.Value) {
			e.Value = `"` + e.Value + `"`
		}
		out = append(out, e)
	}
	return out
}

// NormalizeFile reads a file, normalizes its entries, and writes the result back.
func NormalizeFile(path string, opts ...NormalizeOption) error {
	entries, err := ParseFile(path)
	if err != nil {
		return err
	}
	normalized := Normalize(entries, opts...)
	return WriteFile(path, normalized)
}

func needsQuoting(v string) bool {
	if strings.HasPrefix(v, `"`) {
		return false // already quoted
	}
	for _, r := range v {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

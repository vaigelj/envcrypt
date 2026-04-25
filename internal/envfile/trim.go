package envfile

import (
	"strings"
)

// TrimOption configures the Trim operation.
type TrimOption func(*trimConfig)

type trimConfig struct {
	keys    map[string]bool
	exclude map[string]bool
	cutset  string
}

// WithTrimKeys limits trimming to the specified keys.
func WithTrimKeys(keys ...string) TrimOption {
	return func(c *trimConfig) {
		for _, k := range keys {
			c.keys[k] = true
		}
	}
}

// WithTrimExclude excludes the specified keys from trimming.
func WithTrimExclude(keys ...string) TrimOption {
	return func(c *trimConfig) {
		for _, k := range keys {
			c.exclude[k] = true
		}
	}
}

// WithTrimCutset sets a custom cutset of characters to trim (default: whitespace).
func WithTrimCutset(cutset string) TrimOption {
	return func(c *trimConfig) {
		c.cutset = cutset
	}
}

// Trim removes leading and trailing whitespace (or a custom cutset) from
// entry values. Options control which keys are affected.
func Trim(entries []Entry, opts ...TrimOption) []Entry {
	cfg := &trimConfig{
		keys:    make(map[string]bool),
		exclude: make(map[string]bool),
	}
	for _, o := range opts {
		o(cfg)
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		if cfg.exclude[e.Key] {
			out[i] = e
			continue
		}
		if len(cfg.keys) > 0 && !cfg.keys[e.Key] {
			out[i] = e
			continue
		}
		v := e.Value
		if cfg.cutset != "" {
			v = strings.Trim(v, cfg.cutset)
		} else {
			v = strings.TrimSpace(v)
		}
		out[i] = Entry{Key: e.Key, Value: v, Comment: e.Comment}
	}
	return out
}

// TrimFile reads entries from path, trims values, and writes the result back.
func TrimFile(path string, opts ...TrimOption) error {
	entries, err := ParseFile(path)
	if err != nil {
		return err
	}
	trimmed := Trim(entries, opts...)
	return WriteFile(path, trimmed)
}

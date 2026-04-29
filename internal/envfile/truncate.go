package envfile

import (
	"fmt"
	"strings"
)

// TruncateOption configures how Truncate behaves.
type TruncateOption func(*truncateConfig)

type truncateConfig struct {
	maxLen  int
	suffix  string
	keys    map[string]struct{}
	exclude map[string]struct{}
}

// WithTruncateMaxLen sets the maximum value length (default: 64).
func WithTruncateMaxLen(n int) TruncateOption {
	return func(c *truncateConfig) { c.maxLen = n }
}

// WithTruncateSuffix sets the suffix appended when a value is truncated (default: "...").
func WithTruncateSuffix(s string) TruncateOption {
	return func(c *truncateConfig) { c.suffix = s }
}

// WithTruncateKeys limits truncation to the specified keys only.
func WithTruncateKeys(keys ...string) TruncateOption {
	return func(c *truncateConfig) {
		for _, k := range keys {
			c.keys[k] = struct{}{}
		}
	}
}

// WithTruncateExclude skips the specified keys during truncation.
func WithTruncateExclude(keys ...string) TruncateOption {
	return func(c *truncateConfig) {
		for _, k := range keys {
			c.exclude[k] = struct{}{}
		}
	}
}

// Truncate shortens env entry values that exceed the configured maximum length.
func Truncate(entries []Entry, opts ...TruncateOption) ([]Entry, error) {
	cfg := &truncateConfig{
		maxLen:  64,
		suffix:  "...",
		keys:    make(map[string]struct{}),
		exclude: make(map[string]struct{}),
	}
	for _, o := range opts {
		o(cfg)
	}
	if cfg.maxLen <= 0 {
		return nil, fmt.Errorf("truncate: maxLen must be positive, got %d", cfg.maxLen)
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		if _, skip := cfg.exclude[e.Key]; skip {
			out[i] = e
			continue
		}
		if len(cfg.keys) > 0 {
			if _, ok := cfg.keys[e.Key]; !ok {
				out[i] = e
				continue
			}
		}
		out[i] = e
		if len(e.Value) > cfg.maxLen {
			cutAt := cfg.maxLen - len(cfg.suffix)
			if cutAt < 0 {
				cutAt = 0
			}
			out[i].Value = strings.TrimRight(e.Value[:cutAt], " ") + cfg.suffix
		}
	}
	return out, nil
}

// TruncateFile reads, truncates, and writes an env file in place.
func TruncateFile(path string, opts ...TruncateOption) error {
	entries, err := ParseFile(path)
	if err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	result, err := Truncate(entries, opts...)
	if err != nil {
		return err
	}
	return WriteFile(path, result)
}

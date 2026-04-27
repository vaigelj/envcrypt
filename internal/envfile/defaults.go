package envfile

import (
	"fmt"
	"strings"
)

// DefaultEntry describes a key with its default value and an optional comment.
type DefaultEntry struct {
	Key     string
	Default string
	Comment string
}

// ApplyDefaultsOption configures ApplyDefaults behaviour.
type ApplyDefaultsOption func(*applyDefaultsConfig)

type applyDefaultsConfig struct {
	overwrite bool
}

// WithDefaultsOverwrite causes ApplyDefaults to overwrite existing keys.
func WithDefaultsOverwrite() ApplyDefaultsOption {
	return func(c *applyDefaultsConfig) { c.overwrite = true }
}

// ApplyDefaults merges a list of DefaultEntry values into entries.
// By default it only sets keys that are absent or empty.
// Returns the updated slice and a summary of keys that were applied.
func ApplyDefaults(entries []Entry, defaults []DefaultEntry, opts ...ApplyDefaultsOption) ([]Entry, []string, error) {
	cfg := &applyDefaultsConfig{}
	for _, o := range opts {
		o(cfg)
	}

	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	var applied []string

	for _, d := range defaults {
		if strings.TrimSpace(d.Key) == "" {
			return nil, nil, fmt.Errorf("defaults: empty key is not allowed")
		}
		if idx, exists := index[d.Key]; exists {
			if cfg.overwrite {
				entries[idx].Value = d.Default
				if d.Comment != "" {
					entries[idx].Comment = d.Comment
				}
				applied = append(applied, d.Key)
			}
		} else {
			entries = append(entries, Entry{
				Key:     d.Key,
				Value:   d.Default,
				Comment: d.Comment,
			})
			index[d.Key] = len(entries) - 1
			applied = append(applied, d.Key)
		}
	}

	return entries, applied, nil
}

// ApplyDefaultsFile reads path, applies defaults, and writes the result back.
func ApplyDefaultsFile(path string, defaults []DefaultEntry, opts ...ApplyDefaultsOption) ([]string, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return nil, err
	}
	updated, applied, err := ApplyDefaults(entries, defaults, opts...)
	if err != nil {
		return nil, err
	}
	if err := WriteFile(path, updated); err != nil {
		return nil, err
	}
	return applied, nil
}

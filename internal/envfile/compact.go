package envfile

import (
	"strings"
)

// CompactOption configures the Compact operation.
type CompactOption func(*compactConfig)

type compactConfig struct {
	removeComments bool
	removeBlanks   bool
	dedupeKeys     bool
}

// WithRemoveComments strips comment lines from the output.
func WithRemoveComments() CompactOption {
	return func(c *compactConfig) { c.removeComments = true }
}

// WithRemoveBlanks strips blank lines from the output.
func WithRemoveBlanks() CompactOption {
	return func(c *compactConfig) { c.removeBlanks = true }
}

// WithDedupeKeys removes duplicate keys, keeping the last occurrence.
func WithDedupeKeys() CompactOption {
	return func(c *compactConfig) { c.dedupeKeys = true }
}

// Compact reduces an env file by applying the requested cleanup passes.
// It returns a new slice of Entry values.
func Compact(entries []Entry, opts ...CompactOption) []Entry {
	cfg := &compactConfig{}
	for _, o := range opts {
		o(cfg)
	}

	result := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if cfg.removeComments && strings.HasPrefix(strings.TrimSpace(e.Key), "#") {
			continue
		}
		if cfg.removeBlanks && strings.TrimSpace(e.Key) == "" && strings.TrimSpace(e.Value) == "" {
			continue
		}
		result = append(result, e)
	}

	if cfg.dedupeKeys {
		seen := make(map[string]int)
		for i, e := range result {
			if e.Key != "" {
				seen[e.Key] = i
			}
		}
		filtered := result[:0]
		for i, e := range result {
			if e.Key == "" {
				filtered = append(filtered, e)
				continue
			}
			if seen[e.Key] == i {
				filtered = append(filtered, e)
			}
		}
		result = filtered
	}

	return result
}

// CompactFile reads, compacts, and rewrites an env file in place.
func CompactFile(path string, opts ...CompactOption) error {
	entries, err := ParseFile(path)
	if err != nil {
		return err
	}
	compacted := Compact(entries, opts...)
	return WriteFile(path, compacted)
}

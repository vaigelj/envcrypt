package envfile

import (
	"regexp"
	"strings"
)

// FilterOption configures a Filter operation.
type FilterOption func(*filterConfig)

type filterConfig struct {
	keys      []string
	prefix    string
	suffix    string
	pattern   string
	exclude   bool
}

// WithFilterKeys restricts results to the given keys.
func WithFilterKeys(keys ...string) FilterOption {
	return func(c *filterConfig) { c.keys = keys }
}

// WithFilterPrefix keeps only entries whose key starts with prefix.
func WithFilterPrefix(prefix string) FilterOption {
	return func(c *filterConfig) { c.prefix = prefix }
}

// WithFilterSuffix keeps only entries whose key ends with suffix.
func WithFilterSuffix(suffix string) FilterOption {
	return func(c *filterConfig) { c.suffix = suffix }
}

// WithFilterPattern keeps only entries whose key matches the regex pattern.
func WithFilterPattern(pattern string) FilterOption {
	return func(c *filterConfig) { c.pattern = pattern }
}

// WithFilterExclude inverts the filter, removing matched entries.
func WithFilterExclude() FilterOption {
	return func(c *filterConfig) { c.exclude = true }
}

// Filter returns a subset of entries based on the provided options.
// By default it is an inclusive filter; use WithFilterExclude to invert.
func Filter(entries []Entry, opts ...FilterOption) ([]Entry, error) {
	cfg := &filterConfig{}
	for _, o := range opts {
		o(cfg)
	}

	keySet := make(map[string]bool, len(cfg.keys))
	for _, k := range cfg.keys {
		keySet[k] = true
	}

	var re *regexp.Regexp
	if cfg.pattern != "" {
		var err error
		re, err = regexp.Compile(cfg.pattern)
		if err != nil {
			return nil, err
		}
	}

	var result []Entry
	for _, e := range entries {
		matched := matchesFilter(e.Key, cfg, keySet, re)
		if cfg.exclude {
			matched = !matched
		}
		if matched {
			result = append(result, e)
		}
	}
	return result, nil
}

func matchesFilter(key string, cfg *filterConfig, keySet map[string]bool, re *regexp.Regexp) bool {
	if len(keySet) > 0 && !keySet[key] {
		return false
	}
	if cfg.prefix != "" && !strings.HasPrefix(key, cfg.prefix) {
		return false
	}
	if cfg.suffix != "" && !strings.HasSuffix(key, cfg.suffix) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	return true
}

// FilterFile reads entries from path, applies Filter, and writes results back.
func FilterFile(path string, opts ...FilterOption) error {
	entries, err := ParseFile(path)
	if err != nil {
		return err
	}
	filtered, err := Filter(entries, opts...)
	if err != nil {
		return err
	}
	return WriteFile(path, filtered)
}

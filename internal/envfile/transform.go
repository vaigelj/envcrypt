package envfile

import (
	"strings"
)

// TransformFunc is a function that transforms a value given its key.
type TransformFunc func(key, value string) string

// TransformOptions controls which keys are transformed.
type TransformOptions struct {
	// Keys limits transformation to these keys; empty means all keys.
	Keys []string
	// Exclude skips these keys.
	Exclude []string
}

// Transform applies fn to each entry in pairs according to opts.
func Transform(pairs []Entry, fn TransformFunc, opts TransformOptions) []Entry {
	excludeSet := toSet(opts.Exclude)
	keySet := toSet(opts.Keys)

	result := make([]Entry, len(pairs))
	for i, e := range pairs {
		if excludeSet[e.Key] {
			result[i] = e
			continue
		}
		if len(keySet) > 0 && !keySet[e.Key] {
			result[i] = e
			continue
		}
		result[i] = Entry{Key: e.Key, Value: fn(e.Key, e.Value)}
	}
	return result
}

// UppercaseValues returns a TransformFunc that uppercases values.
func UppercaseValues() TransformFunc {
	return func(_, v string) string { return strings.ToUpper(v) }
}

// TrimValues returns a TransformFunc that trims whitespace from values.
func TrimValues() TransformFunc {
	return func(_, v string) string { return strings.TrimSpace(v) }
}

// PrefixValues returns a TransformFunc that prepends prefix to each value.
func PrefixValues(prefix string) TransformFunc {
	return func(_, v string) string { return prefix + v }
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

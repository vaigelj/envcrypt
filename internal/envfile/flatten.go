package envfile

import (
	"fmt"
	"strings"
)

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	// Separator is placed between key segments (default: "_").
	Separator string
	// Uppercase forces all resulting keys to uppercase.
	Uppercase bool
	// Prefix is prepended to every resulting key.
	Prefix string
}

// Flatten takes a map of dot-notation or slash-notation keys and expands them
// into a flat list of env-style entries. Nested keys like "db.host" become
// "DB_HOST" when Uppercase is true and Separator is "_".
func Flatten(entries []Entry, opts FlattenOptions) []Entry {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		key := flattenKey(e.Key, opts)
		out = append(out, Entry{Key: key, Value: e.Value, Comment: e.Comment})
	}
	return out
}

// FlattenMap converts a nested string map (using dot-separated keys) into a
// flat slice of Entry values suitable for writing to a .env file.
func FlattenMap(m map[string]string, opts FlattenOptions) []Entry {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	out := make([]Entry, 0, len(m))
	for k, v := range m {
		key := flattenKey(k, opts)
		out = append(out, Entry{Key: key, Value: v})
	}
	return out
}

// UnflattenToMap reconstructs a nested map from flat env entries using the
// given separator to split keys into segments.
func UnflattenToMap(entries []Entry, separator string) map[string]interface{} {
	if separator == "" {
		separator = "_"
	}
	root := make(map[string]interface{})
	for _, e := range entries {
		parts := strings.Split(e.Key, separator)
		setNested(root, parts, e.Value)
	}
	return root
}

// FormatFlattened returns a human-readable summary of flattened keys.
func FormatFlattened(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}

func flattenKey(key string, opts FlattenOptions) string {
	// Normalise separators: replace dots and slashes with the target separator.
	r := strings.NewReplacer(".", opts.Separator, "/", opts.Separator)
	key = r.Replace(key)
	if opts.Uppercase {
		key = strings.ToUpper(key)
	}
	if opts.Prefix != "" {
		key = opts.Prefix + key
	}
	return key
}

func setNested(m map[string]interface{}, parts []string, value string) {
	if len(parts) == 1 {
		m[parts[0]] = value
		return
	}
	sub, ok := m[parts[0]]
	if !ok {
		sub = make(map[string]interface{})
		m[parts[0]] = sub
	}
	if subMap, ok := sub.(map[string]interface{}); ok {
		setNested(subMap, parts[1:], value)
	}
}

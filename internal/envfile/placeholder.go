package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// placeholderRe matches patterns like {{VAR_NAME}} or {{ VAR_NAME }}
var placeholderRe = regexp.MustCompile(`\{\{\s*([A-Z_][A-Z0-9_]*)\s*\}\}`)

// ResolvePlaceholders replaces {{KEY}} tokens in entry values using the
// provided env map as a lookup table. Unresolved placeholders are left intact
// unless strict is true, in which case an error is returned.
func ResolvePlaceholders(entries []Entry, env map[string]string, strict bool) ([]Entry, error) {
	// Build a lookup from the entries themselves so self-referencing works.
	lookup := make(map[string]string, len(env))
	for k, v := range env {
		lookup[k] = v
	}
	for _, e := range entries {
		if _, exists := lookup[e.Key]; !exists {
			lookup[e.Key] = e.Value
		}
	}

	resolved := make([]Entry, len(entries))
	for i, e := range entries {
		val, err := resolveSingle(e.Value, lookup, strict)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", e.Key, err)
		}
		resolved[i] = Entry{Key: e.Key, Value: val}
	}
	return resolved, nil
}

// ResolvePlaceholdersString resolves placeholders within a single string.
func ResolvePlaceholdersString(s string, env map[string]string, strict bool) (string, error) {
	return resolveSingle(s, env, strict)
}

func resolveSingle(s string, lookup map[string]string, strict bool) (string, error) {
	var missingKeys []string
	result := placeholderRe.ReplaceAllStringFunc(s, func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-2])
		if val, ok := lookup[key]; ok {
			return val
		}
		missingKeys = append(missingKeys, key)
		return match
	})
	if strict && len(missingKeys) > 0 {
		return "", fmt.Errorf("unresolved placeholders: %s", strings.Join(missingKeys, ", "))
	}
	return result, nil
}

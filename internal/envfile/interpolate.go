package envfile

import (
	"fmt"
	"os"
	"strings"
)

// InterpolateOptions controls interpolation behavior.
type InterpolateOptions struct {
	// Strict causes an error when a referenced variable is not found.
	Strict bool
	// Environ supplements the entries map with OS environment variables as fallback.
	Environ bool
}

// Interpolate replaces ${VAR} and $VAR references in entry values using the
// provided map of entries. If opts.Environ is true, os.Getenv is used as a
// fallback for missing keys. If opts.Strict is true, any unresolved reference
// returns an error.
func Interpolate(entries []Entry, opts InterpolateOptions) ([]Entry, error) {
	// Build lookup map from current entries.
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		resolved, err := interpolateString(e.Value, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", e.Key, err)
		}
		out[i] = Entry{Key: e.Key, Value: resolved}
	}
	return out, nil
}

// InterpolateFile reads a file, interpolates values, and writes the result back.
func InterpolateFile(path string, opts InterpolateOptions) error {
	entries, err := ParseFile(path)
	if err != nil {
		return err
	}
	resolved, err := Interpolate(entries, opts)
	if err != nil {
		return err
	}
	return WriteFile(path, resolved)
}

func interpolateString(s string, lookup map[string]string, opts InterpolateOptions) (string, error) {
	var sb strings.Builder
	i := 0
	for i < len(s) {
		if s[i] != '$' {
			sb.WriteByte(s[i])
			i++
			continue
		}
		// '$' found
		i++
		if i >= len(s) {
			sb.WriteByte('$')
			break
		}
		var name string
		if s[i] == '{' {
			// ${VAR} form
			end := strings.IndexByte(s[i:], '}')
			if end < 0 {
				sb.WriteByte('$')
				continue
			}
			name = s[i+1 : i+end]
			i += end + 1
		} else {
			// $VAR form
			j := i
			for j < len(s) && isVarChar(s[j]) {
				j++
			}
			if j == i {
				sb.WriteByte('$')
				continue
			}
			name = s[i:j]
			i = j
		}
		if val, ok := lookup[name]; ok {
			sb.WriteString(val)
		} else if opts.Environ {
			sb.WriteString(os.Getenv(name))
		} else if opts.Strict {
			return "", fmt.Errorf("undefined variable %q", name)
		} else {
			sb.WriteString("$" + name)
		}
	}
	return sb.String(), nil
}

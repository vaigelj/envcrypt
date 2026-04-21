package envfile

import (
	"fmt"
	"os"
	"strings"
)

// ResolveMode controls how missing keys are handled during resolution.
type ResolveMode int

const (
	// ResolveModeStrict returns an error for any unresolved reference.
	ResolveModeStrict ResolveMode = iota
	// ResolveModeLoose leaves unresolved references as-is.
	ResolveModeLoose
)

// ResolveOptions configures the Resolve operation.
type ResolveOptions struct {
	// Mode controls strict vs loose resolution.
	Mode ResolveMode
	// Environ supplements the entries map with OS environment variables.
	Environ bool
}

// Resolve expands ${VAR} and $VAR references within entry values using
// other entries in the same map. If opts.Environ is true, OS environment
// variables are also consulted as a fallback.
//
// References are resolved in a single pass; circular references are not
// followed and will be left as-is in loose mode or return an error in
// strict mode.
func Resolve(entries []Entry, opts ResolveOptions) ([]Entry, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	resolved := make([]Entry, len(entries))
	for i, e := range entries {
		val, err := resolveValue(e.Value, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", e.Key, err)
		}
		resolved[i] = Entry{Key: e.Key, Value: val}
	}
	return resolved, nil
}

// ResolveFile reads a .env file, resolves variable references, and returns
// the expanded entries.
func ResolveFile(path string, opts ResolveOptions) ([]Entry, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return nil, err
	}
	return Resolve(entries, opts)
}

func resolveValue(val string, lookup map[string]string, opts ResolveOptions) (string, error) {
	var sb strings.Builder
	s := val
	for len(s) > 0 {
		start := strings.Index(s, "$")
		if start == -1 {
			sb.WriteString(s)
			break
		}
		sb.WriteString(s[:start])
		s = s[start:]

		var key string
		var consumed int
		if strings.HasPrefix(s, "${") {
			end := strings.Index(s, "}")
			if end == -1 {
				sb.WriteString(s)
				break
			}
			key = s[2:end]
			consumed = end + 1
		} else {
			i := 1
			for i < len(s) && isVarChar(s[i]) {
				i++
			}
			key = s[1:i]
			consumed = i
		}

		if key == "" {
			sb.WriteByte('$')
			s = s[1:]
			continue
		}

		if v, ok := lookup[key]; ok {
			sb.WriteString(v)
		} else if opts.Environ {
			if v, ok := os.LookupEnv(key); ok {
				sb.WriteString(v)
			} else if opts.Mode == ResolveModeStrict {
				return "", fmt.Errorf("unresolved variable %q", key)
			} else {
				sb.WriteString(s[:consumed])
			}
		} else if opts.Mode == ResolveModeStrict {
			return "", fmt.Errorf("unresolved variable %q", key)
		} else {
			sb.WriteString(s[:consumed])
		}
		s = s[consumed:]
	}
	return sb.String(), nil
}

func isVarChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') || c == '_'
}

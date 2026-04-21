package envfile

import (
	"fmt"
	"os"
	"strings"
)

// InjectOptions controls how environment variables are injected into the process.
type InjectOptions struct {
	// Overwrite existing environment variables if true.
	Overwrite bool
	// Prefix is prepended to every key before injection.
	Prefix string
	// Only inject keys present in this set (nil means all).
	Only map[string]bool
}

// Inject sets environment variables from entries into the current process.
// It respects InjectOptions for overwrite, prefix filtering, and key allowlist.
func Inject(entries []Entry, opts InjectOptions) error {
	for _, e := range entries {
		key := e.Key
		if opts.Only != nil && !opts.Only[key] {
			continue
		}
		if opts.Prefix != "" {
			key = opts.Prefix + key
		}
		if !opts.Overwrite {
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
		}
		if err := os.Setenv(key, e.Value); err != nil {
			return fmt.Errorf("inject: setenv %q: %w", key, err)
		}
	}
	return nil
}

// InjectFile reads a .env file and injects its variables into the process.
func InjectFile(path string, opts InjectOptions) error {
	entries, err := ParseFile(path)
	if err != nil {
		return fmt.Errorf("inject: parse file: %w", err)
	}
	return Inject(entries, opts)
}

// Snapshot of env keys injected — returns a restore function that unsets them.
func InjectWithRollback(entries []Entry, opts InjectOptions) (func(), error) {
	var injected []string
	for _, e := range entries {
		key := e.Key
		if opts.Only != nil && !opts.Only[key] {
			continue
		}
		if opts.Prefix != "" {
			key = opts.Prefix + key
		}
		if !opts.Overwrite {
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
		}
		if err := os.Setenv(key, e.Value); err != nil {
			return nil, fmt.Errorf("inject: setenv %q: %w", key, err)
		}
		injected = append(injected, key)
	}
	rollback := func() {
		for _, k := range injected {
			_ = os.Unsetenv(k)
		}
	}
	return rollback, nil
}

// toSet converts a comma-separated string into a lookup map.
func parseKeySet(csv string) map[string]bool {
	if strings.TrimSpace(csv) == "" {
		return nil
	}
	m := make(map[string]bool)
	for _, k := range strings.Split(csv, ",") {
		m[strings.TrimSpace(k)] = true
	}
	return m
}

// Package envfile provides utilities for parsing, manipulating, and writing .env files.
package envfile

import (
	"fmt"
	"strings"
)

// SetOption configures the behaviour of Set and SetFile.
type SetOption func(*setConfig)

type setConfig struct {
	overwrite bool
	comment   string
}

// WithOverwrite allows Set to replace an existing key's value.
func WithOverwrite() SetOption {
	return func(c *setConfig) { c.overwrite = true }
}

// WithComment attaches an inline comment to the new or updated entry.
// The comment should not include the leading '#'.
func WithComment(comment string) SetOption {
	return func(c *setConfig) { c.comment = comment }
}

// Set adds or updates a single key=value pair in the provided entries slice.
// By default it will not overwrite an existing key; pass WithOverwrite() to
// allow updates. The returned slice preserves the original order and appends
// new entries at the end.
//
// An error is returned if the key is empty or contains whitespace.
func Set(entries []Entry, key, value string, opts ...SetOption) ([]Entry, error) {
	if err := validateKey(key); err != nil {
		return nil, err
	}

	cfg := &setConfig{}
	for _, o := range opts {
		o(cfg)
	}

	for i, e := range entries {
		if e.Key == key {
			if !cfg.overwrite {
				return nil, fmt.Errorf("set: key %q already exists (use WithOverwrite to replace)", key)
			}
			entries[i].Value = value
			if cfg.comment != "" {
				entries[i].Comment = cfg.comment
			}
			return entries, nil
		}
	}

	// Key not found — append.
	entries = append(entries, Entry{
		Key:     key,
		Value:   value,
		Comment: cfg.comment,
	})
	return entries, nil
}

// SetFile reads the .env file at path, applies Set with the given options, and
// writes the result back to the same path.
func SetFile(path, key, value string, opts ...SetOption) error {
	entries, err := ParseFile(path)
	if err != nil {
		return fmt.Errorf("set: read %s: %w", path, err)
	}

	entries, err = Set(entries, key, value, opts...)
	if err != nil {
		return err
	}

	if err := WriteFile(path, entries); err != nil {
		return fmt.Errorf("set: write %s: %w", path, err)
	}
	return nil
}

// Delete removes the entry with the given key from entries.
// It returns the updated slice and true if the key was found, or the original
// slice and false if it was not present.
func Delete(entries []Entry, key string) ([]Entry, bool) {
	for i, e := range entries {
		if e.Key == key {
			return append(entries[:i], entries[i+1:]...), true
		}
	}
	return entries, false
}

// DeleteFile reads the .env file at path, removes the given key, and writes
// the result back. It returns an error if the key does not exist.
func DeleteFile(path, key string) error {
	entries, err := ParseFile(path)
	if err != nil {
		return fmt.Errorf("delete: read %s: %w", path, err)
	}

	updated, found := Delete(entries, key)
	if !found {
		return fmt.Errorf("delete: key %q not found in %s", key, path)
	}

	if err := WriteFile(path, updated); err != nil {
		return fmt.Errorf("delete: write %s: %w", path, err)
	}
	return nil
}

// validateKey returns an error for keys that would produce an invalid .env line.
func validateKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("set: key must not be empty")
	}
	if strings.ContainsAny(key, " \t\n\r=") {
		return fmt.Errorf("set: key %q contains invalid characters", key)
	}
	return nil
}

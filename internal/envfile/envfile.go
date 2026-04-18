// Package envfile provides functionality for parsing, encrypting,
// and decrypting .env files using the envcrypt format.
package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key   string
	Value string
}

// EnvFile holds the parsed entries of a .env file.
type EnvFile struct {
	Entries []Entry
}

// Parse reads a .env file from the given path and returns an EnvFile.
// Lines starting with '#' and empty lines are ignored.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envfile: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("envfile: invalid line %q", line)
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"`)
		entries = append(entries, Entry{Key: key, Value: val})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envfile: scan: %w", err)
	}
	return &EnvFile{Entries: entries}, nil
}

// Write serializes the EnvFile to the given path.
func (e *EnvFile) Write(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("envfile: create %q: %w", path, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, entry := range e.Entries {
		if _, err := fmt.Fprintf(w, "%s=%q\n", entry.Key, entry.Value); err != nil {
			return fmt.Errorf("envfile: write: %w", err)
		}
	}
	return w.Flush()
}

// ToMap converts the EnvFile entries into a map for easy lookup.
func (e *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(e.Entries))
	for _, entry := range e.Entries {
		m[entry.Key] = entry.Value
	}
	return m
}

// Get returns the value for the given key and whether it was found.
func (e *EnvFile) Get(key string) (string, bool) {
	for _, entry := range e.Entries {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return "", false
}

package envfile

import (
	"fmt"
	"os"
)

// ParseFile reads a .env file from disk and returns its key-value pairs.
func ParseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	return Parse(f)
}

// WriteFile writes the given key-value pairs to a .env file at path,
// creating or truncating the file as needed.
func WriteFile(path string, env map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()
	for k, v := range env {
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, v); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	return nil
}

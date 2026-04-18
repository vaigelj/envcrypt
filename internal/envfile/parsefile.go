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

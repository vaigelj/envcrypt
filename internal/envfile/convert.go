package envfile

import (
	"fmt"
	"strings"
)

// ConvertFormat re-encodes entries from one format to another.
// Supported formats: dotenv, json, shell.
func ConvertFormat(entries []Entry, from, to string) (string, error) {
	if from == to {
		return Export(entries, to)
	}

	// Validate source format (entries are already parsed, but we verify the target)
	switch strings.ToLower(to) {
	case "dotenv", "json", "shell":
		// ok
	default:
		return "", fmt.Errorf("unsupported target format %q: must be dotenv, json, or shell", to)
	}

	return Export(entries, to)
}

// ConvertFile reads a file in the given format and writes converted output
// to the destination file in the target format.
func ConvertFile(srcPath, dstPath, from, to string) error {
	entries, err := Import(srcPath, from)
	if err != nil {
		return fmt.Errorf("convert read: %w", err)
	}

	out, err := ConvertFormat(entries, from, to)
	if err != nil {
		return fmt.Errorf("convert format: %w", err)
	}

	if err := writeString(dstPath, out); err != nil {
		return fmt.Errorf("convert write: %w", err)
	}
	return nil
}

// writeString is a small helper that writes a string to a file.
func writeString(path, content string) error {
	return WriteFile(path, parseEntries(content))
}

// parseEntries parses a raw dotenv string into entries.
func parseEntries(raw string) []Entry {
	entries, _ := Parse(strings.NewReader(raw))
	return entries
}

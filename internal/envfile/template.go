package envfile

import (
	"fmt"
	"os"
	"strings"
)

// TemplateEntry represents a single entry in a .env template file.
type TemplateEntry struct {
	Key      string
	Required bool
	Comment  string
}

// ParseTemplate reads a .env.template file where values are either empty or
// contain a description comment. Lines prefixed with "#!" mark required keys.
//
// Example:
//
//	DATABASE_URL=        # required database connection string
//	#! API_KEY=          # required
//	OPTIONAL_FLAG=
func ParseTemplate(path string) ([]TemplateEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("template: read %q: %w", path, err)
	}

	var entries []TemplateEntry
	for i, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		required := false
		if strings.HasPrefix(line, "#!") {
			required = true
			line = strings.TrimSpace(line[2:])
		} else if strings.HasPrefix(line, "#") {
			continue
		}

		// Strip inline comment
		comment := ""
		if idx := strings.Index(line, "#"); idx != -1 {
			comment = strings.TrimSpace(line[idx+1:])
			line = strings.TrimSpace(line[:idx])
		}

		eqIdx := strings.Index(line, "=")
		if eqIdx == -1 {
			return nil, fmt.Errorf("template: line %d: missing '='", i+1)
		}

		key := strings.TrimSpace(line[:eqIdx])
		if key == "" {
			return nil, fmt.Errorf("template: line %d: empty key", i+1)
		}

		entries = append(entries, TemplateEntry{Key: key, Required: required, Comment: comment})
	}
	return entries, nil
}

// CheckTemplate verifies that all required keys from the template exist in env.
// It returns a list of missing required keys.
func CheckTemplate(entries []TemplateEntry, env map[string]string) []string {
	var missing []string
	for _, e := range entries {
		if e.Required {
			if v, ok := env[e.Key]; !ok || strings.TrimSpace(v) == "" {
				missing = append(missing, e.Key)
			}
		}
	}
	return missing
}

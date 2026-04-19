package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ImportFormat represents a supported import format.
type ImportFormat string

const (
	ImportDotenv ImportFormat = "dotenv"
	ImportJSON   ImportFormat = "json"
	ImportShell  ImportFormat = "shell"
)

// Import reads environment variables from src in the given format
// and returns them as a map.
func Import(src, format string) (map[string]string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("import: read %s: %w", src, err)
	}
	switch ImportFormat(format) {
	case ImportDotenv:
		return importDotenv(string(data))
	case ImportJSON:
		return importJSON(data)
	case ImportShell:
		return importShell(string(data))
	default:
		return nil, fmt.Errorf("import: unknown format %q", format)
	}
}

func importDotenv(content string) (map[string]string, error) {
	return Parse(content)
}

func importJSON(data []byte) (map[string]string, error) {
	raw := map[string]interface{}{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("import json: %w", err)
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out, nil
}

func importShell(content string) (map[string]string, error) {
	out := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimPrefix(line, "export ")
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "'\"")
		out[key] = val
	}
	return out, nil
}

package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportFormat defines the output format for exporting env vars.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatShell  ExportFormat = "shell"
)

// Export serializes a map of env vars into the requested format.
func Export(vars map[string]string, format ExportFormat) (string, error) {
	switch format {
	case FormatDotenv:
		return exportDotenv(vars), nil
	case FormatJSON:
		return exportJSON(vars)
	case FormatShell:
		return exportShell(vars), nil
	default:
		return "", fmt.Errorf("unknown export format: %q", format)
	}
}

// ExportFile writes the exported content to a file.
func ExportFile(vars map[string]string, format ExportFormat, path string) error {
	content, err := Export(vars, format)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0600)
}

func exportDotenv(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return sb.String()
}

func exportJSON(vars map[string]string) (string, error) {
	b, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b) + "\n", nil
}

func exportShell(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}

package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var sampleVars = map[string]string{
	"APP_ENV": "production",
	"DB_PASS": "s3cr3t",
}

func TestExportDotenv(t *testing.T) {
	out, err := Export(sampleVars, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range sampleVars {
		expected := k + "=" + v
		if !strings.Contains(out, expected) {
			t.Errorf("expected %q in dotenv output", expected)
		}
	}
}

func TestExportJSON(t *testing.T) {
	out, err := Export(sampleVars, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}
	for k, v := range sampleVars {
		if parsed[k] != v {
			t.Errorf("key %s: got %q want %q", k, parsed[k], v)
		}
	}
}

func TestExportShell(t *testing.T) {
	out, err := Export(sampleVars, FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export ") {
		t.Error("shell export should contain 'export ' prefix")
	}
}

func TestExportUnknownFormat(t *testing.T) {
	_, err := Export(sampleVars, ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestExportFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	if err := ExportFile(sampleVars, FormatDotenv, path); err != nil {
		t.Fatalf("ExportFile error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}
	if len(data) == 0 {
		t.Error("exported file should not be empty")
	}
}

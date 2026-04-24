package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempConvertEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestConvertDotenvToJSON(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	out, err := ConvertFormat(entries, "dotenv", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"FOO"`) || !strings.Contains(out, `"bar"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestConvertDotenvToShell(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
	}
	out, err := ConvertFormat(entries, "dotenv", "shell")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export FOO") {
		t.Errorf("expected shell export, got: %s", out)
	}
}

func TestConvertSameFormat(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}}
	out, err := ConvertFormat(entries, "dotenv", "dotenv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "A=1") {
		t.Errorf("expected dotenv output, got: %s", out)
	}
}

func TestConvertUnsupportedTarget(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}}
	_, err := ConvertFormat(entries, "dotenv", "yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestConvertFile(t *testing.T) {
	src := writeTempConvertEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.json")

	if err := ConvertFile(src, dst, "dotenv", "json"); err != nil {
		t.Fatalf("ConvertFile error: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if !strings.Contains(string(data), "FOO") {
		t.Errorf("expected FOO in output, got: %s", data)
	}
}

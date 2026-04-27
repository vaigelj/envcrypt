package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeDefaultEntries() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestApplyDefaultsAddsAbsent(t *testing.T) {
	entries := makeDefaultEntries()
	defs := []DefaultEntry{
		{Key: "TIMEOUT", Default: "30s", Comment: "request timeout"},
	}
	out, applied, err := ApplyDefaults(entries, defs)
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 1 || applied[0] != "TIMEOUT" {
		t.Fatalf("expected TIMEOUT applied, got %v", applied)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[2].Value != "30s" {
		t.Errorf("expected 30s, got %s", out[2].Value)
	}
}

func TestApplyDefaultsDoesNotOverwriteByDefault(t *testing.T) {
	entries := makeDefaultEntries()
	defs := []DefaultEntry{
		{Key: "HOST", Default: "example.com"},
	}
	out, applied, err := ApplyDefaults(entries, defs)
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 0 {
		t.Fatalf("expected nothing applied, got %v", applied)
	}
	if out[0].Value != "localhost" {
		t.Errorf("HOST should remain localhost, got %s", out[0].Value)
	}
}

func TestApplyDefaultsOverwrite(t *testing.T) {
	entries := makeDefaultEntries()
	defs := []DefaultEntry{
		{Key: "PORT", Default: "9090"},
	}
	out, applied, err := ApplyDefaults(entries, defs, WithDefaultsOverwrite())
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 1 || applied[0] != "PORT" {
		t.Fatalf("expected PORT applied, got %v", applied)
	}
	if out[1].Value != "9090" {
		t.Errorf("expected 9090, got %s", out[1].Value)
	}
}

func TestApplyDefaultsEmptyKeyError(t *testing.T) {
	entries := makeDefaultEntries()
	defs := []DefaultEntry{{Key: "", Default: "val"}}
	_, _, err := ApplyDefaults(entries, defs)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestApplyDefaultsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("HOST=localhost\nPORT=8080\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	defs := []DefaultEntry{
		{Key: "DEBUG", Default: "false"},
	}
	applied, err := ApplyDefaultsFile(path, defs)
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 1 || applied[0] != "DEBUG" {
		t.Fatalf("expected DEBUG applied, got %v", applied)
	}
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range entries {
		if e.Key == "DEBUG" && e.Value == "false" {
			found = true
		}
	}
	if !found {
		t.Error("DEBUG=false not persisted to file")
	}
}

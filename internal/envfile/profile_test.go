package envfile

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestListProfiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{".env.dev", ".env.staging", ".env.prod", "README.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("KEY=val\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	profiles, err := ListProfiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(profiles))
	}
	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name
	}
	sort.Strings(names)
	expected := []string{"dev", "prod", "staging"}
	for i, n := range expected {
		if names[i] != n {
			t.Errorf("expected %q, got %q", n, names[i])
		}
	}
}

func TestLoadProfile(t *testing.T) {
	dir := t.TempDir()
	content := "DB_HOST=localhost\nDB_PORT=5432\n"
	if err := os.WriteFile(filepath.Join(dir, ".env.dev"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	m, err := LoadProfile(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", m["DB_HOST"])
	}
	if m["DB_PORT"] != "5432" {
		t.Errorf("expected 5432, got %q", m["DB_PORT"])
	}
}

func TestLoadProfileMissing(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadProfile(dir, "ghost")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestSaveAndLoadProfile(t *testing.T) {
	dir := t.TempDir()
	data := map[string]string{"APP_ENV": "staging", "LOG_LEVEL": "info"}
	if err := SaveProfile(dir, "staging", data); err != nil {
		t.Fatalf("save error: %v", err)
	}
	m, err := LoadProfile(dir, "staging")
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	for k, v := range data {
		if m[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, m[k])
		}
	}
}

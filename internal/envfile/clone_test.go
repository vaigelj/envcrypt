package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForClone(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestCloneAll(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	got := Clone(entries, CloneOptions{})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Value != "bar" || got[1].Value != "qux" {
		t.Errorf("unexpected values: %+v", got)
	}
}

func TestCloneSubset(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "1"},
		{Key: "BAR", Value: "2"},
		{Key: "BAZ", Value: "3"},
	}
	got := Clone(entries, CloneOptions{Keys: []string{"FOO", "BAZ"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Key != "FOO" || got[1].Key != "BAZ" {
		t.Errorf("unexpected keys: %+v", got)
	}
}

func TestCloneStripValues(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET", Value: "s3cr3t"},
		{Key: "OTHER", Value: "visible"},
	}
	got := Clone(entries, CloneOptions{StripValues: true})
	for _, e := range got {
		if e.Value != "" {
			t.Errorf("expected empty value for %q, got %q", e.Key, e.Value)
		}
	}
}

func TestCloneFileNoOverwrite(t *testing.T) {
	src := writeTempEnvForClone(t, "FOO=bar\n")
	dst := filepath.Join(t.TempDir(), "out.env")
	// Create dst so it already exists
	if err := os.WriteFile(dst, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := CloneFile(src, dst, CloneOptions{Overwrite: false})
	if err == nil {
		t.Fatal("expected error when destination exists and Overwrite=false")
	}
}

func TestCloneFileOverwrite(t *testing.T) {
	src := writeTempEnvForClone(t, "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(t.TempDir(), "out.env")
	if err := os.WriteFile(dst, []byte("old=data\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := CloneFile(src, dst, CloneOptions{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := ParseFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 entries after clone, got %d", len(got))
	}
}

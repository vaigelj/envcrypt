package envfile

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestDedupeKeepFirst(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "A", Value: "99"},
	}
	res := Dedupe(entries, DedupeKeepFirst)
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "1" {
		t.Errorf("expected first A=1, got %s", res.Entries[0].Value)
	}
	if !slices.Contains(res.Duplicates, "A") {
		t.Errorf("expected A in duplicates")
	}
}

func TestDedupeKeepLast(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "A", Value: "99"},
	}
	res := Dedupe(entries, DedupeKeepLast)
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "99" {
		t.Errorf("expected last A=99, got %s", res.Entries[0].Value)
	}
}

func TestDedupeNoDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "X", Value: "a"},
		{Key: "Y", Value: "b"},
	}
	res := Dedupe(entries, DedupeKeepFirst)
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", res.Duplicates)
	}
}

func TestDedupePreservesOrder(t *testing.T) {
	entries := []Entry{
		{Key: "C", Value: "3"},
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "A", Value: "99"},
	}
	res := Dedupe(entries, DedupeKeepFirst)
	got := make([]string, len(res.Entries))
	for i, e := range res.Entries {
		got[i] = e.Key
	}
	want := []string{"C", "A", "B"}
	if !slices.Equal(got, want) {
		t.Errorf("order mismatch: got %v want %v", got, want)
	}
}

func TestDedupeFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "FOO=1\nBAR=2\nFOO=3\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	res, err := DedupeFile(path, DedupeKeepLast)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "3" {
		t.Errorf("expected FOO=3 after keep-last, got %s", res.Entries[0].Value)
	}
	if !slices.Contains(res.Duplicates, "FOO") {
		t.Errorf("expected FOO in duplicates")
	}
}

package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeReorderEntries() []Entry {
	return []Entry{
		{Key: "C", Value: "3"},
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "D", Value: "4"},
	}
}

func TestReorderBasic(t *testing.T) {
	entries := makeReorderEntries()
	out, err := Reorder(entries, []string{"A", "B", "C", "D"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"A", "B", "C", "D"}
	for i, e := range out {
		if e.Key != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Key, want[i])
		}
	}
}

func TestReorderPartialOrderRemainsRelative(t *testing.T) {
	entries := makeReorderEntries()
	out, err := Reorder(entries, []string{"D", "A"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// D and A first, then C and B in original relative order
	want := []string{"D", "A", "C", "B"}
	for i, e := range out {
		if e.Key != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Key, want[i])
		}
	}
}

func TestReorderMissingKeyError(t *testing.T) {
	entries := makeReorderEntries()
	_, err := Reorder(entries, []string{"A", "Z"})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestReorderMissingKeyOk(t *testing.T) {
	entries := makeReorderEntries()
	out, err := Reorder(entries, []string{"Z", "A"}, WithReorderMissingOk())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "A" {
		t.Errorf("expected first key A, got %q", out[0].Key)
	}
}

func TestReorderFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("C=3\nA=1\nB=2\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ReorderFile(path, []string{"A", "B", "C"}); err != nil {
		t.Fatalf("ReorderFile: %v", err)
	}
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"A", "B", "C"}
	for i, e := range entries {
		if e.Key != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Key, want[i])
		}
	}
}

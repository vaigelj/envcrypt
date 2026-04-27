package envfile

import (
	"os"
	"testing"
)

func makeCompactEntries() []Entry {
	return []Entry{
		{Key: "# comment", Value: ""},
		{Key: "", Value: ""},
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
		{Key: "FOO", Value: "overridden"},
		{Key: "", Value: ""},
	}
}

func TestCompactRemoveComments(t *testing.T) {
	out := Compact(makeCompactEntries(), WithRemoveComments())
	for _, e := range out {
		if e.Key == "# comment" {
			t.Fatal("expected comment to be removed")
		}
	}
	if len(out) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(out))
	}
}

func TestCompactRemoveBlanks(t *testing.T) {
	out := Compact(makeCompactEntries(), WithRemoveBlanks())
	for _, e := range out {
		if e.Key == "" && e.Value == "" {
			t.Fatal("expected blank lines to be removed")
		}
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
}

func TestCompactDedupeKeys(t *testing.T) {
	out := Compact(makeCompactEntries(), WithDedupeKeys())
	for _, e := range out {
		if e.Key == "FOO" && e.Value != "overridden" {
			t.Fatalf("expected last FOO value 'overridden', got %q", e.Value)
		}
	}
	count := 0
	for _, e := range out {
		if e.Key == "FOO" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 FOO entry, got %d", count)
	}
}

func TestCompactCombined(t *testing.T) {
	out := Compact(makeCompactEntries(),
		WithRemoveComments(),
		WithRemoveBlanks(),
		WithDedupeKeys(),
	)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries after full compact, got %d", len(out))
	}
}

func TestCompactFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "compact*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("# comment\nFOO=bar\nFOO=baz\n\nBAR=qux\n")
	f.Close()

	if err := CompactFile(f.Name(), WithRemoveComments(), WithRemoveBlanks(), WithDedupeKeys()); err != nil {
		t.Fatal(err)
	}
	entries, err := ParseFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

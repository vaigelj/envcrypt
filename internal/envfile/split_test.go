package envfile

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func makeSplitEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envcrypt"},
		{Key: "APP_ENV", Value: "prod"},
		{Key: "STANDALONE", Value: "yes"},
	}
}

func TestSplitByPrefix(t *testing.T) {
	groups := Split(makeSplitEntries())
	if len(groups["DB"]) != 2 {
		t.Fatalf("expected 2 DB entries, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 2 {
		t.Fatalf("expected 2 APP entries, got %d", len(groups["APP"]))
	}
	if len(groups["_default"]) != 1 {
		t.Fatalf("expected 1 _default entry, got %d", len(groups["_default"]))
	}
}

func TestSplitCustomSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "DB.HOST", Value: "localhost"},
		{Key: "DB.PORT", Value: "5432"},
		{Key: "APP.NAME", Value: "envcrypt"},
	}
	groups := Split(entries, WithSplitSeparator("."))
	if len(groups["DB"]) != 2 {
		t.Fatalf("expected 2 DB entries, got %d", len(groups["DB"]))
	}
}

func TestSplitNoSeparatorGoesToDefault(t *testing.T) {
	entries := []Entry{{Key: "NOPREFIX", Value: "val"}}
	groups := Split(entries)
	if len(groups["_default"]) != 1 {
		t.Fatal("expected entry in _default group")
	}
}

func TestSplitFile(t *testing.T) {
	src := filepath.Join(t.TempDir(), "input.env")
	if err := os.WriteFile(src, []byte("DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=envcrypt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	written, err := SplitFile(src, out)
	if err != nil {
		t.Fatalf("SplitFile: %v", err)
	}
	sort.Strings(written)
	if len(written) != 2 {
		t.Fatalf("expected 2 files, got %d", len(written))
	}
}

func TestSplitFileNoOverwrite(t *testing.T) {
	src := filepath.Join(t.TempDir(), "input.env")
	_ = os.WriteFile(src, []byte("DB_HOST=localhost\n"), 0o644)
	out := t.TempDir()
	_ = os.WriteFile(filepath.Join(out, "DB.env"), []byte("existing"), 0o644)
	_, err := SplitFile(src, out)
	if err == nil {
		t.Fatal("expected error when file exists without --overwrite")
	}
}

func TestSplitFileOverwrite(t *testing.T) {
	src := filepath.Join(t.TempDir(), "input.env")
	_ = os.WriteFile(src, []byte("DB_HOST=localhost\n"), 0o644)
	out := t.TempDir()
	_ = os.WriteFile(filepath.Join(out, "DB.env"), []byte("old"), 0o644)
	_, err := SplitFile(src, out, WithSplitOverwrite())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

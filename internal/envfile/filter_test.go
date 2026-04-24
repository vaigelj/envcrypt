package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeFilterEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestFilterByKeys(t *testing.T) {
	entries := makeFilterEntries()
	result, err := Filter(entries, WithFilterKeys("DB_HOST", "APP_NAME"))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilterByPrefix(t *testing.T) {
	entries := makeFilterEntries()
	result, err := Filter(entries, WithFilterPrefix("DB_"))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key != "DB_HOST" && e.Key != "DB_PORT" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestFilterBySuffix(t *testing.T) {
	entries := makeFilterEntries()
	result, err := Filter(entries, WithFilterSuffix("_ENV"))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].Key != "APP_ENV" {
		t.Fatalf("expected APP_ENV, got %+v", result)
	}
}

func TestFilterByPattern(t *testing.T) {
	entries := makeFilterEntries()
	result, err := Filter(entries, WithFilterPattern(`^APP_`))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilterInvalidPattern(t *testing.T) {
	entries := makeFilterEntries()
	_, err := Filter(entries, WithFilterPattern(`[invalid`))
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestFilterExclude(t *testing.T) {
	entries := makeFilterEntries()
	result, err := Filter(entries, WithFilterPrefix("DB_"), WithFilterExclude())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key == "DB_HOST" || e.Key == "DB_PORT" {
			t.Errorf("excluded key %q still present", e.Key)
		}
	}
}

func TestFilterFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("DB_HOST=localhost\nAPP_NAME=myapp\nSECRET_KEY=abc\n"), 0o600)

	if err := FilterFile(path, WithFilterPrefix("DB_")); err != nil {
		t.Fatal(err)
	}
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Key != "DB_HOST" {
		t.Fatalf("expected only DB_HOST, got %+v", entries)
	}
}

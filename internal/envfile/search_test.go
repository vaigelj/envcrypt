package envfile

import (
	"os"
	"testing"
)

func writeTempEnvForSearch(t *testing.T, content string) string {
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

var searchEntries = []Entry{
	{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
	{Key: "SECRET_KEY", Value: "supersecret"},
	{Key: "APP_PORT", Value: "8080"},
	{Key: "APP_HOST", Value: "localhost"},
}

func TestSearchByKeySubstring(t *testing.T) {
	res, err := Search(searchEntries, "test.env", SearchOptions{KeyPattern: "APP_"})
	if err != nil || len(res) != 2 {
		t.Fatalf("expected 2 results, got %d (err=%v)", len(res), err)
	}
}

func TestSearchByValueSubstring(t *testing.T) {
	res, err := Search(searchEntries, "test.env", SearchOptions{ValuePattern: "localhost"})
	if err != nil || len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	res, err := Search(searchEntries, "test.env", SearchOptions{KeyPattern: "app_port", CaseSensitive: false})
	if err != nil || len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
}

func TestSearchRegex(t *testing.T) {
	res, err := Search(searchEntries, "test.env", SearchOptions{KeyPattern: "^APP_", UseRegex: true})
	if err != nil || len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestSearchNoMatch(t *testing.T) {
	res, err := Search(searchEntries, "test.env", SearchOptions{KeyPattern: "NONEXISTENT"})
	if err != nil || len(res) != 0 {
		t.Fatalf("expected 0 results, got %d", len(res))
	}
}

func TestSearchFile(t *testing.T) {
	path := writeTempEnvForSearch(t, "FOO=bar\nBAZ=qux\n")
	res, err := SearchFile(path, SearchOptions{KeyPattern: "FOO"})
	if err != nil || len(res) != 1 || res[0].Key != "FOO" {
		t.Fatalf("unexpected result: %v %v", res, err)
	}
}

func TestSearchInvalidRegex(t *testing.T) {
	_, err := Search(searchEntries, "test.env", SearchOptions{KeyPattern: "[", UseRegex: true})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

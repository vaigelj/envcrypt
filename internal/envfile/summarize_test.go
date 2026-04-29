package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForSummarize(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestSummarizeBasic(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "envcrypt"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "EMPTY_VAL", Value: ""},
		{Key: "# a comment", Comment: true},
	}
	r := Summarize(entries)
	if r.TotalKeys != 3 {
		t.Errorf("TotalKeys = %d, want 3", r.TotalKeys)
	}
	if r.CommentLines != 1 {
		t.Errorf("CommentLines = %d, want 1", r.CommentLines)
	}
	if r.EmptyValues != 1 {
		t.Errorf("EmptyValues = %d, want 1", r.EmptyValues)
	}
	if len(r.SensitiveKeys) == 0 {
		t.Error("expected DB_PASSWORD to be detected as sensitive")
	}
}

func TestSummarizeDuplicates(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "1"},
		{Key: "FOO", Value: "2"},
		{Key: "BAR", Value: "3"},
	}
	r := Summarize(entries)
	if r.UniqueKeys != 2 {
		t.Errorf("UniqueKeys = %d, want 2", r.UniqueKeys)
	}
	if len(r.DuplicateKeys) != 1 || r.DuplicateKeys[0] != "FOO" {
		t.Errorf("DuplicateKeys = %v, want [FOO]", r.DuplicateKeys)
	}
}

func TestSummarizeLongestKey(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "x"},
		{Key: "LONG_KEY_NAME", Value: "y"},
		{Key: "MED", Value: "z"},
	}
	r := Summarize(entries)
	if r.LongestKey != "LONG_KEY_NAME" {
		t.Errorf("LongestKey = %q, want LONG_KEY_NAME", r.LongestKey)
	}
}

func TestSummarizeFile(t *testing.T) {
	content := "APP=hello\nSECRET_KEY=abc\n# comment\nEMPTY=\n"
	p := writeTempEnvForSummarize(t, content)
	r, err := SummarizeFile(p)
	if err != nil {
		t.Fatalf("SummarizeFile error: %v", err)
	}
	if r.TotalKeys != 3 {
		t.Errorf("TotalKeys = %d, want 3", r.TotalKeys)
	}
	if r.EmptyValues != 1 {
		t.Errorf("EmptyValues = %d, want 1", r.EmptyValues)
	}
}

func TestSummarizeString(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "FOO", Value: "baz"},
	}
	r := Summarize(entries)
	s := r.String()
	if !strings.Contains(s, "Duplicate keys") {
		t.Errorf("String() missing duplicate keys section: %s", s)
	}
}

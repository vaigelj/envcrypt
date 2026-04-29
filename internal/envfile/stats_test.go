package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeStatsEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "supersecretvalue123"},
		{Key: "EMPTY_VAR", Value: ""},
		{Key: "APP_NAME", Value: "duplicate"},
		{Key: "", Value: "# a comment"},
	}
}

func TestStatsTotal(t *testing.T) {
	s := GatherStats(makeStatsEntries())
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
}

func TestStatsWithValuesAndEmpty(t *testing.T) {
	s := GatherStats(makeStatsEntries())
	if s.WithValues != 4 {
		t.Errorf("expected WithValues=4, got %d", s.WithValues)
	}
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
}

func TestStatsComments(t *testing.T) {
	s := GatherStats(makeStatsEntries())
	if s.Comments != 1 {
		t.Errorf("expected Comments=1, got %d", s.Comments)
	}
}

func TestStatsUniqueAndDuplicates(t *testing.T) {
	s := GatherStats(makeStatsEntries())
	// APP_NAME appears twice → 1 duplicate; others are unique
	if s.Duplicates != 1 {
		t.Errorf("expected Duplicates=1, got %d", s.Duplicates)
	}
	if s.Unique != 3 {
		t.Errorf("expected Unique=3, got %d", s.Unique)
	}
}

func TestStatsLongestKey(t *testing.T) {
	s := GatherStats(makeStatsEntries())
	if s.LongestKey != "SECRET_KEY" {
		t.Errorf("expected LongestKey=SECRET_KEY, got %s", s.LongestKey)
	}
}

func TestStatsTopKeys(t *testing.T) {
	s := GatherStats(makeStatsEntries())
	if len(s.TopKeys) == 0 {
		t.Fatal("expected at least one top key")
	}
	// SECRET_KEY has the longest value among unique keys
	if s.TopKeys[0] != "SECRET_KEY" {
		t.Errorf("expected first TopKey=SECRET_KEY, got %s", s.TopKeys[0])
	}
}

func TestGatherStatsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "APP=hello\nDB=world\nEMPTY=\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	s, err := GatherStatsFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
}

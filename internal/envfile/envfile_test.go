package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcrypt/internal/envfile"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestParseBasic(t *testing.T) {
	p := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	ef, err := envfile.Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "DB_HOST" || ef.Entries[0].Value != "localhost" {
		t.Errorf("unexpected entry: %+v", ef.Entries[0])
	}
}

func TestParseIgnoresCommentsAndBlanks(t *testing.T) {
	p := writeTempEnv(t, "# comment\n\nFOO=bar\n")
	ef, err := envfile.Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ef.Entries))
	}
}

func TestParseInvalidLine(t *testing.T) {
	p := writeTempEnv(t, "INVALID_LINE\n")
	_, err := envfile.Parse(p)
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}

func TestWriteRoundtrip(t *testing.T) {
	p := writeTempEnv(t, "KEY1=value1\nKEY2=value2\n")
	ef, err := envfile.Parse(p)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	out := filepath.Join(t.TempDir(), ".env.out")
	if err := ef.Write(out); err != nil {
		t.Fatalf("write: %v", err)
	}

	ef2, err := envfile.Parse(out)
	if err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	if len(ef2.Entries) != len(ef.Entries) {
		t.Fatalf("entry count mismatch: %d vs %d", len(ef2.Entries), len(ef.Entries))
	}
	for i, e := range ef.Entries {
		if ef2.Entries[i].Key != e.Key || ef2.Entries[i].Value != e.Value {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, ef2.Entries[i], e)
		}
	}
}

func TestToMap(t *testing.T) {
	p := writeTempEnv(t, "A=1\nB=2\n")
	ef, _ := envfile.Parse(p)
	m := ef.ToMap()
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
}

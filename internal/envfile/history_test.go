package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForHistory(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestAppendAndLoadHistory(t *testing.T) {
	path := writeTempEnvForHistory(t)
	vals := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := AppendHistory(path, "initial", vals); err != nil {
		t.Fatal(err)
	}
	hf, err := LoadHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(hf.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(hf.Entries))
	}
	if hf.Entries[0].Label != "initial" {
		t.Errorf("expected label 'initial', got %q", hf.Entries[0].Label)
	}
	if hf.Entries[0].Values["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
}

func TestAppendHistoryMultiple(t *testing.T) {
	path := writeTempEnvForHistory(t)
	for i, label := range []string{"v1", "v2", "v3"} {
		_ = i
		if err := AppendHistory(path, label, map[string]string{"KEY": label}); err != nil {
			t.Fatal(err)
		}
	}
	hf, err := LoadHistory(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(hf.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(hf.Entries))
	}
	if hf.Entries[2].Label != "v3" {
		t.Errorf("expected last label v3")
	}
}

func TestClearHistory(t *testing.T) {
	path := writeTempEnvForHistory(t)
	_ = AppendHistory(path, "x", map[string]string{})
	if err := ClearHistory(path); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadHistory(path); err == nil {
		t.Error("expected error after clear")
	}
}

func TestLoadHistoryMissing(t *testing.T) {
	_, err := LoadHistory("/nonexistent/path.env")
	if err == nil {
		t.Error("expected error for missing history")
	}
}

func TestHistoryDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.env")
	_ = os.WriteFile(path, []byte("A=1"), 0644)
	_ = AppendHistory(path, "init", map[string]string{"A": "1"})
	files, err := HistoryDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 history file, got %d", len(files))
	}
}

package envfile

import (
	"os"
	"testing"
)

func writeTempEnvForRename(t *testing.T, content string) string {
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

func TestRenameKeyBasic(t *testing.T) {
	m := map[string]string{"OLD": "value", "OTHER": "x"}
	if err := RenameKey(m, "OLD", "NEW", false); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["OLD"]; ok {
		t.Error("old key should be removed")
	}
	if m["NEW"] != "value" {
		t.Errorf("expected NEW=value, got %q", m["NEW"])
	}
}

func TestRenameKeyMissing(t *testing.T) {
	m := map[string]string{"A": "1"}
	if err := RenameKey(m, "MISSING", "NEW", false); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestRenameKeyConflict(t *testing.T) {
	m := map[string]string{"OLD": "1", "NEW": "2"}
	if err := RenameKey(m, "OLD", "NEW", false); err == nil {
		t.Error("expected conflict error")
	}
	if err := RenameKey(m, "OLD", "NEW", true); err != nil {
		t.Errorf("unexpected error with overwrite: %v", err)
	}
}

func TestRenameFile(t *testing.T) {
	path := writeTempEnvForRename(t, "FOO=bar\nBAZ=qux\n")
	if err := RenameFile(path, "FOO", "FOO_RENAMED", false); err != nil {
		t.Fatal(err)
	}
	m, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["FOO"]; ok {
		t.Error("old key should not exist")
	}
	if m["FOO_RENAMED"] != "bar" {
		t.Errorf("expected FOO_RENAMED=bar, got %q", m["FOO_RENAMED"])
	}
}

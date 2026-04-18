package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForSnapshot(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestTakeSnapshot(t *testing.T) {
	path := writeTempEnvForSnapshot(t, "FOO=bar\nBAZ=qux\n")
	snap, err := TakeSnapshot(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Entries["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", snap.Entries["FOO"])
	}
	if snap.Source != path {
		t.Errorf("expected source %s, got %s", path, snap.Source)
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	path := writeTempEnvForSnapshot(t, "KEY=value\n")
	snap, err := TakeSnapshot(path)
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(t.TempDir(), "snap.json")
	if err := SaveSnapshot(snap, dest); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadSnapshot(dest)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Entries["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", loaded.Entries["KEY"])
	}
}

func TestDiffSnapshot(t *testing.T) {
	old := &Snapshot{Entries: map[string]string{"A": "1", "B": "2"}}
	new := &Snapshot{Entries: map[string]string{"A": "changed", "C": "3"}}
	changes := DiffSnapshot(old, new)
	if len(changes) == 0 {
		t.Fatal("expected changes, got none")
	}
	kinds := map[ChangeKind]int{}
	for _, c := range changes {
		kinds[c.Kind]++
	}
	if kinds[Updated] != 1 {
		t.Errorf("expected 1 updated, got %d", kinds[Updated])
	}
	if kinds[Added] != 1 {
		t.Errorf("expected 1 added, got %d", kinds[Added])
	}
	if kinds[Removed] != 1 {
		t.Errorf("expected 1 removed, got %d", kinds[Removed])
	}
}

func TestLoadSnapshotMissing(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

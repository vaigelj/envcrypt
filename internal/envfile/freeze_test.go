package envfile

import (
	"testing"
)

func TestFreezeAndLoad(t *testing.T) {
	dir := t.TempDir()
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}

	f, err := Freeze(dir, "prod", "prod.env", entries)
	if err != nil {
		t.Fatalf("Freeze: %v", err)
	}
	if f.Source != "prod.env" {
		t.Errorf("source = %q, want %q", f.Source, "prod.env")
	}
	if f.Checksum == "" {
		t.Error("expected non-empty checksum")
	}

	loaded, err := LoadFrozen(dir, "prod")
	if err != nil {
		t.Fatalf("LoadFrozen: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Errorf("entries len = %d, want 2", len(loaded.Entries))
	}
	if loaded.Checksum != f.Checksum {
		t.Errorf("checksum mismatch")
	}
}

func TestFreezeIsTampered(t *testing.T) {
	dir := t.TempDir()
	entries := []Entry{{Key: "SECRET", Value: "abc"}}

	f, err := Freeze(dir, "staging", "staging.env", entries)
	if err != nil {
		t.Fatalf("Freeze: %v", err)
	}
	if f.IsTampered() {
		t.Error("expected not tampered")
	}

	// Mutate entries after freeze
	f.Entries[0].Value = "changed"
	if !f.IsTampered() {
		t.Error("expected tampered after mutation")
	}
}

func TestLoadFrozenMissing(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadFrozen(dir, "nonexistent")
	if err == nil {
		t.Error("expected error for missing frozen file")
	}
}

func TestDeleteFrozen(t *testing.T) {
	dir := t.TempDir()
	entries := []Entry{{Key: "X", Value: "1"}}
	_, err := Freeze(dir, "tmp", "tmp.env", entries)
	if err != nil {
		t.Fatalf("Freeze: %v", err)
	}
	if err := DeleteFrozen(dir, "tmp"); err != nil {
		t.Fatalf("DeleteFrozen: %v", err)
	}
	// Deleting again should not error
	if err := DeleteFrozen(dir, "tmp"); err != nil {
		t.Errorf("second DeleteFrozen: %v", err)
	}
}

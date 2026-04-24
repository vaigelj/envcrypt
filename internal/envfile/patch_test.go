package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func basePatchEntries() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestPatchSet(t *testing.T) {
	out, err := Patch(basePatchEntries(), []PatchInstruction{
		{Op: PatchSet, Key: "PORT", Value: "9999"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out[1].Value != "9999" {
		t.Errorf("expected PORT=9999, got %s", out[1].Value)
	}
}

func TestPatchSetNewKey(t *testing.T) {
	out, err := Patch(basePatchEntries(), []PatchInstruction{
		{Op: PatchSet, Key: "NEW_VAR", Value: "hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
	if out[3].Key != "NEW_VAR" || out[3].Value != "hello" {
		t.Errorf("unexpected last entry: %+v", out[3])
	}
}

func TestPatchDelete(t *testing.T) {
	out, err := Patch(basePatchEntries(), []PatchInstruction{
		{Op: PatchDelete, Key: "DEBUG"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Key == "DEBUG" {
			t.Error("DEBUG should have been deleted")
		}
	}
}

func TestPatchRename(t *testing.T) {
	out, err := Patch(basePatchEntries(), []PatchInstruction{
		{Op: PatchRename, Key: "HOST", NewKey: "DB_HOST"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out[0].Key)
	}
}

func TestPatchRenameMissingKey(t *testing.T) {
	_, err := Patch(basePatchEntries(), []PatchInstruction{
		{Op: PatchRename, Key: "NOPE", NewKey: "OTHER"},
	})
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestPatchUnknownOp(t *testing.T) {
	_, err := Patch(basePatchEntries(), []PatchInstruction{
		{Op: PatchOp("upsert"), Key: "X"},
	})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatchFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("HOST=localhost\nPORT=5432\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := PatchFile(path, []PatchInstruction{
		{Op: PatchSet, Key: "PORT", Value: "8080"},
		{Op: PatchDelete, Key: "HOST"},
	})
	if err != nil {
		t.Fatal(err)
	}
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Key != "PORT" || entries[0].Value != "8080" {
		t.Errorf("unexpected entries after patch: %+v", entries)
	}
}

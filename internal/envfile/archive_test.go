package envfile

import (
	"os"
	"testing"
)

func TestSaveAndLoadArchive(t *testing.T) {
	dir := t.TempDir()
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	if err := SaveArchive(dir, "release-1", entries); err != nil {
		t.Fatalf("SaveArchive: %v", err)
	}
	a, err := LoadArchive(dir, "release-1")
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if a.Name != "release-1" {
		t.Errorf("name: got %q want %q", a.Name, "release-1")
	}
	if len(a.Entries) != 2 {
		t.Errorf("entries len: got %d want 2", len(a.Entries))
	}
	if a.Entries[0].Key != "FOO" || a.Entries[0].Value != "bar" {
		t.Errorf("unexpected entry: %+v", a.Entries[0])
	}
}

func TestLoadArchiveMissing(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadArchive(dir, "nope")
	if err == nil {
		t.Fatal("expected error for missing archive")
	}
}

func TestListArchives(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"v2", "v1", "v3"} {
		if err := SaveArchive(dir, name, nil); err != nil {
			t.Fatalf("SaveArchive %s: %v", name, err)
		}
	}
	names, err := ListArchives(dir)
	if err != nil {
		t.Fatalf("ListArchives: %v", err)
	}
	if len(names) != 3 {
		t.Fatalf("expected 3 archives, got %d", len(names))
	}
	if names[0] != "v1" || names[1] != "v2" || names[2] != "v3" {
		t.Errorf("unexpected order: %v", names)
	}
}

func TestListArchivesEmpty(t *testing.T) {
	dir := t.TempDir()
	names, err := ListArchives(dir)
	if err != nil {
		t.Fatalf("ListArchives: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestDeleteArchive(t *testing.T) {
	dir := t.TempDir()
	if err := SaveArchive(dir, "snap", nil); err != nil {
		t.Fatalf("SaveArchive: %v", err)
	}
	if err := DeleteArchive(dir, "snap"); err != nil {
		t.Fatalf("DeleteArchive: %v", err)
	}
	_, err := LoadArchive(dir, "snap")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDeleteArchiveMissing(t *testing.T) {
	dir := t.TempDir()
	if err := DeleteArchive(dir, "ghost"); err == nil {
		t.Fatal("expected error deleting non-existent archive")
	}
}

func TestSaveArchiveEmptyName(t *testing.T) {
	dir := t.TempDir()
	if err := SaveArchive(dir, "", nil); err == nil {
		t.Fatal("expected error for empty archive name")
	}
}

func TestArchiveTimestamp(t *testing.T) {
	dir := t.TempDir()
	if err := SaveArchive(dir, "ts-test", nil); err != nil {
		t.Fatalf("SaveArchive: %v", err)
	}
	a, err := LoadArchive(dir, "ts-test")
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if a.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt timestamp")
	}
	_ = os.RemoveAll(dir)
}

package envfile

import (
	"testing"
	"time"
)

func TestCreateAndLoadBackup(t *testing.T) {
	dir := t.TempDir()
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	b, err := CreateBackup(dir, entries, "initial")
	if err != nil {
		t.Fatalf("CreateBackup: %v", err)
	}
	if b.Label != "initial" {
		t.Errorf("expected label 'initial', got %q", b.Label)
	}
	if len(b.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(b.Entries))
	}

	loaded, err := LoadBackup(dir, b.ID)
	if err != nil {
		t.Fatalf("LoadBackup: %v", err)
	}
	if loaded.ID != b.ID {
		t.Errorf("ID mismatch: got %q want %q", loaded.ID, b.ID)
	}
	if loaded.Entries[0].Key != "FOO" {
		t.Errorf("expected FOO, got %q", loaded.Entries[0].Key)
	}
}

func TestListBackups(t *testing.T) {
	dir := t.TempDir()
	entries := []Entry{{Key: "A", Value: "1"}}

	_, err := CreateBackup(dir, entries, "first")
	if err != nil {
		t.Fatalf("CreateBackup first: %v", err)
	}
	time.Sleep(2 * time.Millisecond)
	_, err = CreateBackup(dir, entries, "second")
	if err != nil {
		t.Fatalf("CreateBackup second: %v", err)
	}

	list, err := ListBackups(dir)
	if err != nil {
		t.Fatalf("ListBackups: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 backups, got %d", len(list))
	}
	// newest first
	if list[0].Label != "second" {
		t.Errorf("expected newest first, got label %q", list[0].Label)
	}
}

func TestListBackupsEmpty(t *testing.T) {
	dir := t.TempDir()
	list, err := ListBackups(dir)
	if err != nil {
		t.Fatalf("ListBackups: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

func TestDeleteBackup(t *testing.T) {
	dir := t.TempDir()
	entries := []Entry{{Key: "X", Value: "y"}}
	b, err := CreateBackup(dir, entries, "to-delete")
	if err != nil {
		t.Fatalf("CreateBackup: %v", err)
	}
	if err := DeleteBackup(dir, b.ID); err != nil {
		t.Fatalf("DeleteBackup: %v", err)
	}
	list, _ := ListBackups(dir)
	if len(list) != 0 {
		t.Errorf("expected 0 backups after delete, got %d", len(list))
	}
}

func TestLoadBackupMissing(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadBackup(dir, "nonexistent")
	if err == nil {
		t.Error("expected error for missing backup")
	}
}

func TestDeleteBackupMissing(t *testing.T) {
	dir := t.TempDir()
	err := DeleteBackup(dir, "ghost")
	if err == nil {
		t.Error("expected error for missing backup")
	}
}

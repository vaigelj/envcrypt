package envfile

import (
	"os"
	"testing"
)

func TestAddAndGetGroup(t *testing.T) {
	dir := t.TempDir()
	if err := AddGroup(dir, "db", []string{"DB_HOST", "DB_PORT"}); err != nil {
		t.Fatalf("AddGroup: %v", err)
	}
	g, err := GetGroup(dir, "db")
	if err != nil {
		t.Fatalf("GetGroup: %v", err)
	}
	if len(g.Keys) != 2 || g.Keys[0] != "DB_HOST" {
		t.Errorf("unexpected keys: %v", g.Keys)
	}
}

func TestAddGroupOverwrites(t *testing.T) {
	dir := t.TempDir()
	_ = AddGroup(dir, "db", []string{"DB_HOST"})
	_ = AddGroup(dir, "db", []string{"DB_HOST", "DB_PASS"})
	g, _ := GetGroup(dir, "db")
	if len(g.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(g.Keys))
	}
}

func TestRemoveGroup(t *testing.T) {
	dir := t.TempDir()
	_ = AddGroup(dir, "db", []string{"DB_HOST"})
	if err := RemoveGroup(dir, "db"); err != nil {
		t.Fatalf("RemoveGroup: %v", err)
	}
	_, err := GetGroup(dir, "db")
	if err == nil {
		t.Error("expected error after removal")
	}
}

func TestRemoveGroupMissing(t *testing.T) {
	dir := t.TempDir()
	if err := RemoveGroup(dir, "nope"); err == nil {
		t.Error("expected error for missing group")
	}
}

func TestLoadGroupsMissing(t *testing.T) {
	dir := t.TempDir()
	groups, err := LoadGroups(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected empty, got %v", groups)
	}
}

func TestGroupsSorted(t *testing.T) {
	dir := t.TempDir()
	_ = AddGroup(dir, "zebra", []string{"Z"})
	_ = AddGroup(dir, "alpha", []string{"A"})
	groups, _ := LoadGroups(dir)
	if groups[0].Name != "alpha" {
		t.Errorf("expected sorted groups, got %v", groups[0].Name)
	}
}

func TestGroupPersistence(t *testing.T) {
	dir := t.TempDir()
	_ = AddGroup(dir, "cache", []string{"REDIS_URL"})
	// reload from disk
	groups, err := LoadGroups(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(groups) != 1 || groups[0].Name != "cache" {
		t.Errorf("unexpected groups after reload: %v", groups)
	}
	_ = os.Remove(groupsPath(dir))
}

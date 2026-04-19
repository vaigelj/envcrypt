package envfile

import (
	"os"
	"testing"
)

func TestAddAndGetTag(t *testing.T) {
	ts := &TagStore{}
	ts.AddTag("secrets", []string{"API_KEY", "DB_PASS"})
	tag := ts.GetTag("secrets")
	if tag == nil {
		t.Fatal("expected tag")
	}
	if len(tag.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(tag.Keys))
	}
}

func TestAddTagOverwrites(t *testing.T) {
	ts := &TagStore{}
	ts.AddTag("group", []string{"A"})
	ts.AddTag("group", []string{"B", "C"})
	if len(ts.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(ts.Tags))
	}
	if ts.Tags[0].Keys[0] != "B" {
		t.Fatal("expected overwritten keys")
	}
}

func TestRemoveTag(t *testing.T) {
	ts := &TagStore{}
	ts.AddTag("x", []string{"K"})
	if !ts.RemoveTag("x") {
		t.Fatal("expected true")
	}
	if ts.GetTag("x") != nil {
		t.Fatal("expected nil after removal")
	}
	if ts.RemoveTag("x") {
		t.Fatal("expected false for missing tag")
	}
}

func TestSaveAndLoadTags(t *testing.T) {
	dir := t.TempDir()
	ts := &TagStore{}
	ts.AddTag("infra", []string{"HOST", "PORT"})
	if err := SaveTags(dir, ts); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadTags(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(loaded.Tags))
	}
	if loaded.Tags[0].Name != "infra" {
		t.Fatal("wrong tag name")
	}
}

func TestLoadTagsMissing(t *testing.T) {
	dir := t.TempDir()
	ts, err := LoadTags(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ts.Tags) != 0 {
		t.Fatal("expected empty store")
	}
}

func TestLoadTagsCorrupt(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(tagsPath(dir), []byte("not json"), 0600)
	_, err := LoadTags(dir)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

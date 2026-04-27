package envfile

import (
	"os"
	"testing"
)

func TestSetAndGetAnnotation(t *testing.T) {
	dir := t.TempDir()
	ann := Annotation{
		Description: "Database URL",
		Owner:       "platform-team",
		Deprecated:  false,
		Tags:        []string{"db", "infra"},
	}
	if err := SetAnnotation(dir, "DATABASE_URL", ann); err != nil {
		t.Fatalf("SetAnnotation: %v", err)
	}
	got, ok, err := GetAnnotation(dir, "DATABASE_URL")
	if err != nil {
		t.Fatalf("GetAnnotation: %v", err)
	}
	if !ok {
		t.Fatal("expected annotation to exist")
	}
	if got.Description != ann.Description {
		t.Errorf("description: got %q want %q", got.Description, ann.Description)
	}
	if got.Owner != ann.Owner {
		t.Errorf("owner: got %q want %q", got.Owner, ann.Owner)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "db" {
		t.Errorf("tags mismatch: %v", got.Tags)
	}
}

func TestGetAnnotationMissing(t *testing.T) {
	dir := t.TempDir()
	_, ok, err := GetAnnotation(dir, "MISSING_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected annotation to be absent")
	}
}

func TestRemoveAnnotation(t *testing.T) {
	dir := t.TempDir()
	if err := SetAnnotation(dir, "API_KEY", Annotation{Description: "secret"}); err != nil {
		t.Fatal(err)
	}
	if err := RemoveAnnotation(dir, "API_KEY"); err != nil {
		t.Fatal(err)
	}
	_, ok, err := GetAnnotation(dir, "API_KEY")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected annotation to be removed")
	}
}

func TestLoadAnnotationsMissing(t *testing.T) {
	dir := t.TempDir()
	a, err := LoadAnnotations(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a) != 0 {
		t.Errorf("expected empty annotations, got %d", len(a))
	}
}

func TestAnnotationPersistence(t *testing.T) {
	dir := t.TempDir()
	keys := []string{"FOO", "BAR", "BAZ"}
	for _, k := range keys {
		if err := SetAnnotation(dir, k, Annotation{Description: "desc of " + k}); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}
	a, err := LoadAnnotations(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 3 {
		t.Errorf("expected 3 annotations, got %d", len(a))
	}
	_ = os.Remove(annotationsPath(dir))
}

package envfile

import (
	"os"
	"testing"
)

func TestSaveAndLoadPin(t *testing.T) {
	dir := t.TempDir()
	vals := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := SavePin(dir, "v1", vals); err != nil {
		t.Fatalf("SavePin: %v", err)
	}
	pin, err := LoadPin(dir, "v1")
	if err != nil {
		t.Fatalf("LoadPin: %v", err)
	}
	if pin == nil {
		t.Fatal("expected pin, got nil")
	}
	if pin.Values["FOO"] != "bar" {
		t.Errorf("expected bar, got %s", pin.Values["FOO"])
	}
	if pin.Name != "v1" {
		t.Errorf("expected name v1, got %s", pin.Name)
	}
}

func TestLoadPinMissing(t *testing.T) {
	dir := t.TempDir()
	pin, err := LoadPin(dir, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pin != nil {
		t.Fatal("expected nil for missing pin")
	}
}

func TestListPins(t *testing.T) {
	dir := t.TempDir()
	_ = SavePin(dir, "alpha", map[string]string{"A": "1"})
	_ = SavePin(dir, "beta", map[string]string{"B": "2"})
	names, err := ListPins(dir)
	if err != nil {
		t.Fatalf("ListPins: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 pins, got %d", len(names))
	}
}

func TestListPinsEmpty(t *testing.T) {
	dir := t.TempDir()
	names, err := ListPins(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected 0 pins, got %d", len(names))
	}
}

func TestDeletePin(t *testing.T) {
	dir := t.TempDir()
	_ = SavePin(dir, "tmp", map[string]string{"X": "y"})
	if err := DeletePin(dir, "tmp"); err != nil {
		t.Fatalf("DeletePin: %v", err)
	}
	pin, _ := LoadPin(dir, "tmp")
	if pin != nil {
		t.Fatal("expected pin to be deleted")
	}
}

func TestDeletePinMissing(t *testing.T) {
	dir := t.TempDir()
	err := DeletePin(dir, "nope")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

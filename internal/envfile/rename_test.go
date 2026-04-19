package envfile

import (
	"os"
	"testing"
)

func TestRenameKeyBasic(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value", "OTHER": "x"}
	res, err := RenameKey(env, "OLD_KEY", "NEW_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Renamed {
		t.Error("expected Renamed=true")
	}
	if _, ok := env["OLD_KEY"]; ok {
		t.Error("old key should be removed")
	}
	if env["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", env["NEW_KEY"])
	}
}

func TestRenameKeyMissing(t *testing.T) {
	env := map[string]string{"A": "1"}
	_, err := RenameKey(env, "MISSING", "B")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenameKeyConflict(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	_, err := RenameKey(env, "A", "B")
	if err == nil {
		t.Fatal("expected error when new key already exists")
	}
}

func TestRenameFile(t *testing.T) {
	f, err := os.CreateTemp("", "rename*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("OLD_NAME=hello\nOTHER=world\n")
	f.Close()
	defer os.Remove(f.Name())

	res, err := RenameFile(f.Name(), "OLD_NAME", "NEW_NAME")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OldKey != "OLD_NAME" || res.NewKey != "NEW_NAME" {
		t.Errorf("unexpected result: %+v", res)
	}

	env, err := ParseFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if env["NEW_NAME"] != "hello" {
		t.Errorf("expected NEW_NAME=hello, got %q", env["NEW_NAME"])
	}
	if _, ok := env["OLD_NAME"]; ok {
		t.Error("OLD_NAME should be gone")
	}
}

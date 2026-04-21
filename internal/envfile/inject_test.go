package envfile

import (
	"os"
	"testing"
)

func TestInjectBasic(t *testing.T) {
	entries := []Entry{
		{Key: "INJECT_FOO", Value: "bar"},
		{Key: "INJECT_BAZ", Value: "qux"},
	}
	t.Cleanup(func() {
		os.Unsetenv("INJECT_FOO")
		os.Unsetenv("INJECT_BAZ")
	})
	if err := Inject(entries, InjectOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("INJECT_FOO"); got != "bar" {
		t.Errorf("INJECT_FOO = %q, want %q", got, "bar")
	}
	if got := os.Getenv("INJECT_BAZ"); got != "qux" {
		t.Errorf("INJECT_BAZ = %q, want %q", got, "qux")
	}
}

func TestInjectNoOverwrite(t *testing.T) {
	os.Setenv("INJECT_EXISTING", "original")
	t.Cleanup(func() { os.Unsetenv("INJECT_EXISTING") })

	entries := []Entry{{Key: "INJECT_EXISTING", Value: "new"}}
	if err := Inject(entries, InjectOptions{Overwrite: false}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("INJECT_EXISTING"); got != "original" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestInjectOverwrite(t *testing.T) {
	os.Setenv("INJECT_OW", "old")
	t.Cleanup(func() { os.Unsetenv("INJECT_OW") })

	entries := []Entry{{Key: "INJECT_OW", Value: "new"}}
	if err := Inject(entries, InjectOptions{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("INJECT_OW"); got != "new" {
		t.Errorf("expected new value, got %q", got)
	}
}

func TestInjectWithPrefix(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("APP_KEY1") })
	entries := []Entry{{Key: "KEY1", Value: "val1"}}
	if err := Inject(entries, InjectOptions{Prefix: "APP_"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("APP_KEY1"); got != "val1" {
		t.Errorf("APP_KEY1 = %q, want %q", got, "val1")
	}
}

func TestInjectOnlySubset(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("INJECT_A")
		os.Unsetenv("INJECT_B")
	})
	entries := []Entry{
		{Key: "INJECT_A", Value: "1"},
		{Key: "INJECT_B", Value: "2"},
	}
	opts := InjectOptions{Only: map[string]bool{"INJECT_A": true}}
	if err := Inject(entries, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("INJECT_A"); got != "1" {
		t.Errorf("INJECT_A = %q, want %q", got, "1")
	}
	if got := os.Getenv("INJECT_B"); got != "" {
		t.Errorf("INJECT_B should not be set, got %q", got)
	}
}

func TestInjectWithRollback(t *testing.T) {
	entries := []Entry{{Key: "INJECT_RB", Value: "temp"}}
	rollback, err := InjectWithRollback(entries, InjectOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("INJECT_RB"); got != "temp" {
		t.Errorf("INJECT_RB = %q, want %q", got, "temp")
	}
	rollback()
	if got := os.Getenv("INJECT_RB"); got != "" {
		t.Errorf("after rollback INJECT_RB should be unset, got %q", got)
	}
}

func TestParseKeySet(t *testing.T) {
	m := parseKeySet("FOO, BAR,BAZ")
	for _, k := range []string{"FOO", "BAR", "BAZ"} {
		if !m[k] {
			t.Errorf("expected key %q in set", k)
		}
	}
	if parseKeySet("") != nil {
		t.Error("expected nil for empty string")
	}
}

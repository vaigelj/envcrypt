package envfile

import (
	"os"
	"testing"
)

func TestTrimAll(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "  hello  "},
		{Key: "BAR", Value: "\t world\t"},
		{Key: "BAZ", Value: "clean"},
	}
	got := Trim(entries)
	expect := []string{"hello", "world", "clean"}
	for i, e := range got {
		if e.Value != expect[i] {
			t.Errorf("entry %d: got %q, want %q", i, e.Value, expect[i])
		}
	}
}

func TestTrimSpecificKeys(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "  hello  "},
		{Key: "BAR", Value: "  world  "},
	}
	got := Trim(entries, WithTrimKeys("FOO"))
	if got[0].Value != "hello" {
		t.Errorf("FOO: got %q, want %q", got[0].Value, "hello")
	}
	if got[1].Value != "  world  " {
		t.Errorf("BAR should be unchanged, got %q", got[1].Value)
	}
}

func TestTrimExclude(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "  hello  "},
		{Key: "BAR", Value: "  world  "},
	}
	got := Trim(entries, WithTrimExclude("BAR"))
	if got[0].Value != "hello" {
		t.Errorf("FOO: got %q, want %q", got[0].Value, "hello")
	}
	if got[1].Value != "  world  " {
		t.Errorf("BAR should be unchanged, got %q", got[1].Value)
	}
}

func TestTrimCutset(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "***secret***"},
	}
	got := Trim(entries, WithTrimCutset("*"))
	if got[0].Value != "secret" {
		t.Errorf("got %q, want %q", got[0].Value, "secret")
	}
}

func TestTrimPreservesComment(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "  val  ", Comment: "keep me"},
	}
	got := Trim(entries)
	if got[0].Comment != "keep me" {
		t.Errorf("comment lost: got %q", got[0].Comment)
	}
}

func TestTrimFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("FOO=  hello  \nBAR=  world  \n")
	f.Close()

	if err := TrimFile(f.Name()); err != nil {
		t.Fatalf("TrimFile: %v", err)
	}
	entries, err := ParseFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Value != strings.TrimSpace(e.Value) {
			t.Errorf("%s value not trimmed: %q", e.Key, e.Value)
		}
	}
}

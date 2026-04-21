package envfile

import (
	"strings"
	"testing"
)

func TestMaskValueFull(t *testing.T) {
	opts := defaultMaskOptions
	opts.Mode = MaskFull
	got := MaskValue("secret123", opts)
	if got != "*********" {
		t.Errorf("expected all asterisks, got %q", got)
	}
}

func TestMaskValuePartial(t *testing.T) {
	opts := MaskOptions{Mode: MaskPartial, MaskChar: '*', RevealLen: 2}
	got := MaskValue("mysecret", opts)
	// expect: my****et
	if !strings.HasPrefix(got, "my") || !strings.HasSuffix(got, "et") {
		t.Errorf("unexpected partial mask: %q", got)
	}
	if len(got) != len("mysecret") {
		t.Errorf("length mismatch: got %d want %d", len(got), len("mysecret"))
	}
}

func TestMaskValuePartialShort(t *testing.T) {
	opts := MaskOptions{Mode: MaskPartial, MaskChar: '*', RevealLen: 3}
	got := MaskValue("ab", opts)
	if got != "**" {
		t.Errorf("short value should be fully masked, got %q", got)
	}
}

func TestMaskValueLength(t *testing.T) {
	opts := MaskOptions{Mode: MaskLength, MaskChar: '#', FixedLen: 6}
	got := MaskValue("anyvalue", opts)
	if got != "######" {
		t.Errorf("expected 6 hashes, got %q", got)
	}
}

func TestMaskValueEmpty(t *testing.T) {
	got := MaskValue("", defaultMaskOptions)
	if got != "" {
		t.Errorf("empty value should remain empty, got %q", got)
	}
}

func TestMaskEntriesSensitiveAuto(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "abc123"},
		{Key: "HOST", Value: "localhost"},
	}
	masked := MaskEntries(entries, defaultMaskOptions)
	if masked[0].Value == "abc123" {
		t.Error("API_KEY should have been masked")
	}
	if masked[1].Value != "localhost" {
		t.Errorf("HOST should not be masked, got %q", masked[1].Value)
	}
}

func TestMaskEntriesSpecificKeys(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
	}
	opts := defaultMaskOptions
	opts.Keys = []string{"HOST"}
	masked := MaskEntries(entries, opts)
	if masked[0].Value == "localhost" {
		t.Error("HOST should have been masked")
	}
	if masked[1].Value != "5432" {
		t.Errorf("PORT should not be masked, got %q", masked[1].Value)
	}
}

func TestMaskSummary(t *testing.T) {
	original := []Entry{
		{Key: "SECRET", Value: "hunter2"},
		{Key: "HOST", Value: "localhost"},
	}
	masked := []Entry{
		{Key: "SECRET", Value: "*******"},
		{Key: "HOST", Value: "localhost"},
	}
	summary := MaskSummary(original, masked)
	if !strings.Contains(summary, "SECRET") {
		t.Errorf("summary should mention SECRET, got: %q", summary)
	}
	if strings.Contains(summary, "HOST") {
		t.Errorf("summary should not mention HOST, got: %q", summary)
	}
}

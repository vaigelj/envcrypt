package envfile

import (
	"testing"
)

func makeSanitizeEntries() []Entry {
	return []Entry{
		{Key: "A", Value: "hello\x00world"},
		{Key: "B", Value: "line1\r\nline2"},
		{Key: "C", Value: "\x01hidden\x02"},
		{Key: "D", Value: "\"quoted\""},
		{Key: "E", Value: "'single'"},
		{Key: "F", Value: "normal"},
	}
}

func TestSanitizeRemoveNullBytes(t *testing.T) {
	entries := makeSanitizeEntries()
	out := Sanitize(entries, WithRemoveNullBytes())
	if out[0].Value != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", out[0].Value)
	}
	// other entries unchanged
	if out[1].Value != entries[1].Value {
		t.Errorf("unexpected change to entry B")
	}
}

func TestSanitizeNormalizeNewlines(t *testing.T) {
	entries := makeSanitizeEntries()
	out := Sanitize(entries, WithNormalizeNewlines())
	if out[1].Value != "line1\nline2" {
		t.Errorf("expected normalized newline, got %q", out[1].Value)
	}
}

func TestSanitizeStripControlChars(t *testing.T) {
	entries := makeSanitizeEntries()
	out := Sanitize(entries, WithStripControlChars())
	if out[2].Value != "hidden" {
		t.Errorf("expected 'hidden', got %q", out[2].Value)
	}
	// tab and newline should be preserved
	tabEntry := []Entry{{Key: "T", Value: "a\tb"}}
	tabOut := Sanitize(tabEntry, WithStripControlChars())
	if tabOut[0].Value != "a\tb" {
		t.Errorf("tab should be preserved, got %q", tabOut[0].Value)
	}
}

func TestSanitizeTrimDoubleQuotes(t *testing.T) {
	entries := makeSanitizeEntries()
	out := Sanitize(entries, WithTrimQuotes())
	if out[3].Value != "quoted" {
		t.Errorf("expected 'quoted', got %q", out[3].Value)
	}
}

func TestSanitizeTrimSingleQuotes(t *testing.T) {
	entries := makeSanitizeEntries()
	out := Sanitize(entries, WithTrimQuotes())
	if out[4].Value != "single" {
		t.Errorf("expected 'single', got %q", out[4].Value)
	}
}

func TestSanitizeCombined(t *testing.T) {
	entries := []Entry{
		{Key: "X", Value: "\"hello\x00\r\nworld\""},
	}
	out := Sanitize(entries,
		WithRemoveNullBytes(),
		WithNormalizeNewlines(),
		WithTrimQuotes(),
	)
	if out[0].Value != "hello\nworld" {
		t.Errorf("expected combined sanitize, got %q", out[0].Value)
	}
}

func TestSanitizePreservesNormal(t *testing.T) {
	entries := makeSanitizeEntries()
	out := Sanitize(entries, WithRemoveNullBytes(), WithStripControlChars())
	if out[5].Value != "normal" {
		t.Errorf("expected 'normal' unchanged, got %q", out[5].Value)
	}
}

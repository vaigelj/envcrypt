package envfile

import (
	"os"
	"testing"
)

func makeTruncateEntries() []Entry {
	return []Entry{
		{Key: "SHORT", Value: "hi"},
		{Key: "LONG_KEY", Value: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789XY"},
		{Key: "EXACT", Value: "exactly64charsXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"},
		{Key: "COMMENT", Value: "normal", Comment: "# keep this"},
	}
}

func TestTruncateShortValuesUnchanged(t *testing.T) {
	out, err := Truncate(makeTruncateEntries())
	if err != nil {
		t.Fatal(err)
	}
	if out[0].Value != "hi" {
		t.Errorf("expected 'hi', got %q", out[0].Value)
	}
}

func TestTruncateLongValue(t *testing.T) {
	out, err := Truncate(makeTruncateEntries(), WithTruncateMaxLen(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(out[1].Value) > 10 {
		t.Errorf("value not truncated: %q", out[1].Value)
	}
	if out[1].Value[len(out[1].Value)-3:] != "..." {
		t.Errorf("expected suffix '...', got %q", out[1].Value)
	}
}

func TestTruncateCustomSuffix(t *testing.T) {
	out, err := Truncate(makeTruncateEntries(), WithTruncateMaxLen(8), WithTruncateSuffix("~"))
	if err != nil {
		t.Fatal(err)
	}
	if out[1].Value[len(out[1].Value)-1:] != "~" {
		t.Errorf("expected suffix '~', got %q", out[1].Value)
	}
}

func TestTruncateSpecificKeys(t *testing.T) {
	out, err := Truncate(makeTruncateEntries(), WithTruncateMaxLen(5), WithTruncateKeys("LONG_KEY"))
	if err != nil {
		t.Fatal(err)
	}
	// SHORT should be unchanged even though it's short
	if out[0].Value != "hi" {
		t.Errorf("SHORT changed unexpectedly: %q", out[0].Value)
	}
	if len(out[1].Value) > 5 {
		t.Errorf("LONG_KEY not truncated: %q", out[1].Value)
	}
}

func TestTruncateExclude(t *testing.T) {
	out, err := Truncate(makeTruncateEntries(), WithTruncateMaxLen(5), WithTruncateExclude("LONG_KEY"))
	if err != nil {
		t.Fatal(err)
	}
	original := makeTruncateEntries()
	if out[1].Value != original[1].Value {
		t.Errorf("excluded key was modified: %q", out[1].Value)
	}
}

func TestTruncateInvalidMaxLen(t *testing.T) {
	_, err := Truncate(makeTruncateEntries(), WithTruncateMaxLen(0))
	if err == nil {
		t.Error("expected error for maxLen=0")
	}
}

func TestTruncateFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	entries := []Entry{
		{Key: "LONG", Value: "aaaaabbbbbcccccdddddeeeeefffff"},
	}
	if err := WriteFile(f.Name(), entries); err != nil {
		t.Fatal(err)
	}
	if err := TruncateFile(f.Name(), WithTruncateMaxLen(10)); err != nil {
		t.Fatal(err)
	}
	result, err := ParseFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(result[0].Value) > 10 {
		t.Errorf("file value not truncated: %q", result[0].Value)
	}
}

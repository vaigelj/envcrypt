package envfile

import (
	"os"
	"testing"
)

func TestParseFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	_, _ = tmp.WriteString("FOO=bar\nBAZ=qux\n")
	tmp.Close()

	m, err := ParseFile(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestParseFileMissing(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInterpolateBasic(t *testing.T) {
	entries := []Entry{
		{Key: "BASE", Value: "/home/user"},
		{Key: "CONF", Value: "${BASE}/.config"},
	}
	out, err := Interpolate(entries, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "/home/user/.config" {
		t.Errorf("expected /home/user/.config, got %q", out[1].Value)
	}
}

func TestInterpolateBareVar(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "URL", Value: "http://$HOST:8080"},
	}
	out, err := Interpolate(entries, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "http://localhost:8080" {
		t.Errorf("got %q", out[1].Value)
	}
}

func TestInterpolateStrictMissingError(t *testing.T) {
	entries := []Entry{
		{Key: "URL", Value: "http://${MISSING_HOST}:8080"},
	}
	_, err := Interpolate(entries, InterpolateOptions{Strict: true})
	if err == nil {
		t.Fatal("expected error for missing variable in strict mode")
	}
}

func TestInterpolateLooseMissingLeaves(t *testing.T) {
	entries := []Entry{
		{Key: "URL", Value: "http://$MISSING:8080"},
	}
	out, err := Interpolate(entries, InterpolateOptions{Strict: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "http://$MISSING:8080" {
		t.Errorf("expected original value, got %q", out[0].Value)
	}
}

func TestInterpolateEnvironFallback(t *testing.T) {
	t.Setenv("ENVCRYPT_TEST_HOST", "envhost")
	entries := []Entry{
		{Key: "URL", Value: "http://${ENVCRYPT_TEST_HOST}:9000"},
	}
	out, err := Interpolate(entries, InterpolateOptions{Environ: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "http://envhost:9000" {
		t.Errorf("got %q", out[0].Value)
	}
}

func TestInterpolateFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "BASE=/opt\nPATH_VAR=${BASE}/bin\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := InterpolateFile(path, InterpolateOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if entries[1].Value != "/opt/bin" {
		t.Errorf("expected /opt/bin, got %q", entries[1].Value)
	}
}

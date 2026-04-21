package envfile

import (
	"os"
	"testing"
)

func TestResolveBraceRef(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "DSN", Value: "postgres://${HOST}/db"},
	}
	out, err := Resolve(entries, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "postgres://localhost/db" {
		t.Errorf("got %q, want %q", out[1].Value, "postgres://localhost/db")
	}
}

func TestResolveBareRef(t *testing.T) {
	entries := []Entry{
		{Key: "PORT", Value: "5432"},
		{Key: "ADDR", Value: "host:$PORT"},
	}
	out, err := Resolve(entries, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "host:5432" {
		t.Errorf("got %q, want %q", out[1].Value, "host:5432")
	}
}

func TestResolveStrictMissingReturnsError(t *testing.T) {
	entries := []Entry{
		{Key: "URL", Value: "http://${MISSING}/path"},
	}
	_, err := Resolve(entries, ResolveOptions{Mode: ResolveModeStrict})
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestResolveLooseMissingLeavesAsIs(t *testing.T) {
	entries := []Entry{
		{Key: "URL", Value: "http://${MISSING}/path"},
	}
	out, err := Resolve(entries, ResolveOptions{Mode: ResolveModeLoose})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "http://${MISSING}/path" {
		t.Errorf("got %q", out[0].Value)
	}
}

func TestResolveEnvironFallback(t *testing.T) {
	os.Setenv("ENVCRYPT_TEST_VAR", "world")
	defer os.Unsetenv("ENVCRYPT_TEST_VAR")

	entries := []Entry{
		{Key: "GREETING", Value: "hello ${ENVCRYPT_TEST_VAR}"},
	}
	out, err := Resolve(entries, ResolveOptions{Environ: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "hello world" {
		t.Errorf("got %q", out[0].Value)
	}
}

func TestResolveNoReferences(t *testing.T) {
	entries := []Entry{
		{Key: "PLAIN", Value: "no refs here"},
	}
	out, err := Resolve(entries, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "no refs here" {
		t.Errorf("got %q", out[0].Value)
	}
}

func TestResolveFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("BASE=hello\nFULL=${BASE}_world\n")
	f.Close()

	out, err := ResolveFile(f.Name(), ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "hello_world" {
		t.Errorf("got %q", out[1].Value)
	}
}

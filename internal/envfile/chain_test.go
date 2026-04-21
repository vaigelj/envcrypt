package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForChain(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestChainResolveEmpty(t *testing.T) {
	c := NewChain()
	entries, err := c.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected no entries, got %d", len(entries))
	}
}

func TestChainResolveSingle(t *testing.T) {
	path := writeTempEnvForChain(t, "FOO=bar\nBAZ=qux\n")
	c := NewChain(path)
	m, err := c.ResolveMap()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestChainResolvePrecedence(t *testing.T) {
	base := writeTempEnvForChain(t, "FOO=base\nSHARED=base\n")
	override := writeTempEnvForChain(t, "SHARED=override\nEXTRA=new\n")
	c := NewChain(base, override)
	m, err := c.ResolveMap()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["FOO"] != "base" {
		t.Errorf("FOO: want %q, got %q", "base", m["FOO"])
	}
	if m["SHARED"] != "override" {
		t.Errorf("SHARED: want %q, got %q", "override", m["SHARED"])
	}
	if m["EXTRA"] != "new" {
		t.Errorf("EXTRA: want %q, got %q", "new", m["EXTRA"])
	}
}

func TestChainResolveThreeLayers(t *testing.T) {
	a := writeTempEnvForChain(t, "A=1\nB=1\nC=1\n")
	b := writeTempEnvForChain(t, "B=2\nC=2\n")
	c := writeTempEnvForChain(t, "C=3\n")
	chain := NewChain(a, b, c)
	m, err := chain.ResolveMap()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["A"] != "1" || m["B"] != "2" || m["C"] != "3" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestChainResolveMissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nonexistent.env")
	c := NewChain(missing)
	_, err := c.Resolve()
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestChainSources(t *testing.T) {
	c := NewChain("a.env", "b.env", "c.env")
	src := c.Sources()
	if len(src) != 3 || src[0] != "a.env" || src[2] != "c.env" {
		t.Errorf("unexpected sources: %v", src)
	}
}

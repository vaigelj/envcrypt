package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCopyTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCopyEnvNoOverwrite(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"A": "99", "C": "3"}
	n := CopyEnv(dst, src, CopyOptions{Overwrite: false})
	if n != 1 {
		t.Fatalf("expected 1 copied, got %d", n)
	}
	if dst["A"] != "1" {
		t.Error("existing key should not be overwritten")
	}
	if dst["C"] != "3" {
		t.Error("new key should be added")
	}
}

func TestCopyEnvOverwrite(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "99", "B": "2"}
	n := CopyEnv(dst, src, CopyOptions{Overwrite: true})
	if n != 2 {
		t.Fatalf("expected 2 copied, got %d", n)
	}
	if dst["A"] != "99" {
		t.Error("key should be overwritten")
	}
}

func TestCopyEnvExclude(t *testing.T) {
	dst := map[string]string{}
	src := map[string]string{"A": "1", "SECRET": "s"}
	n := CopyEnv(dst, src, CopyOptions{Exclude: []string{"SECRET"}})
	if n != 1 {
		t.Fatalf("expected 1 copied, got %d", n)
	}
	if _, ok := dst["SECRET"]; ok {
		t.Error("excluded key should not be copied")
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	srcPath := writeCopyTempEnv(t, dir, "src.env", "FOO=bar\nBAZ=qux\n")
	dstPath := writeCopyTempEnv(t, dir, "dst.env", "FOO=original\n")

	n, err := CopyFile(dstPath, srcPath, CopyOptions{Overwrite: false})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("expected 1 new key copied, got %d", n)
	}

	result, err := ParseFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}
	if result["FOO"] != "original" {
		t.Error("FOO should not be overwritten")
	}
	if result["BAZ"] != "qux" {
		t.Error("BAZ should be copied")
	}
}

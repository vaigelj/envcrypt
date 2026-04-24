package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPromoteAddsNewKeys(t *testing.T) {
	src := []Entry{{Key: "NEW", Value: "1"}, {Key: "ALSO_NEW", Value: "2"}}
	dst := []Entry{{Key: "EXISTING", Value: "x"}}

	out, res, err := Promote(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestPromoteConflictNoOverwrite(t *testing.T) {
	src := []Entry{{Key: "KEY", Value: "new"}}
	dst := []Entry{{Key: "KEY", Value: "old"}}

	_, res, err := Promote(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Conflict) != 1 || res.Conflict[0] != "KEY" {
		t.Fatalf("expected conflict on KEY, got %v", res.Conflict)
	}
	if len(res.Promoted) != 0 {
		t.Fatalf("expected 0 promoted, got %d", len(res.Promoted))
	}
}

func TestPromoteConflictWithOverwrite(t *testing.T) {
	src := []Entry{{Key: "KEY", Value: "new"}}
	dst := []Entry{{Key: "KEY", Value: "old"}}

	out, res, err := Promote(src, dst, WithPromoteOverwrite())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Promoted) != 1 {
		t.Fatalf("expected 1 promoted, got %d", len(res.Promoted))
	}
	if out[0].Value != "new" {
		t.Fatalf("expected value 'new', got %q", out[0].Value)
	}
}

func TestPromoteExclude(t *testing.T) {
	src := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	dst := []Entry{}

	_, res, err := Promote(src, dst, WithPromoteExclude("B"))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "A" {
		t.Fatalf("expected A promoted, got %v", res.Promoted)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "B" {
		t.Fatalf("expected B skipped, got %v", res.Skipped)
	}
}

func TestPromoteFile(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "src.env")
	dstPath := filepath.Join(dir, "dst.env")

	_ = os.WriteFile(srcPath, []byte("NEW_KEY=hello\n"), 0o600)
	_ = os.WriteFile(dstPath, []byte("EXISTING=world\n"), 0o600)

	res, err := PromoteFile(srcPath, dstPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Promoted) != 1 {
		t.Fatalf("expected 1 promoted, got %d", len(res.Promoted))
	}

	entries, err := ParseFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries in dst, got %d", len(entries))
	}
}

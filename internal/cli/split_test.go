package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForSplit(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func runSplitCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envcrypt"}
	root.AddCommand(NewSplitCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"split"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestSplitCmdBasic(t *testing.T) {
	src := writeTempEnvForSplit(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=envcrypt\n")
	out := t.TempDir()
	output, err := runSplitCmd(t, src, "--out", out)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, output)
	}
	if !strings.Contains(output, "wrote 2 file(s)") {
		t.Errorf("expected '2 file(s)' in output, got: %s", output)
	}
	if _, err := os.Stat(filepath.Join(out, "DB.env")); err != nil {
		t.Error("DB.env not created")
	}
	if _, err := os.Stat(filepath.Join(out, "APP.env")); err != nil {
		t.Error("APP.env not created")
	}
}

func TestSplitCmdOverwrite(t *testing.T) {
	src := writeTempEnvForSplit(t, "DB_HOST=localhost\n")
	out := t.TempDir()
	_ = os.WriteFile(filepath.Join(out, "DB.env"), []byte("old"), 0o644)
	_, err := runSplitCmd(t, src, "--out", out)
	if err == nil {
		t.Fatal("expected error without --overwrite")
	}
	_, err = runSplitCmd(t, src, "--out", out, "--overwrite")
	if err != nil {
		t.Fatalf("unexpected error with --overwrite: %v", err)
	}
}

func TestSplitCmdCustomSep(t *testing.T) {
	src := writeTempEnvForSplit(t, "DB.HOST=localhost\nAPP.NAME=envcrypt\n")
	out := t.TempDir()
	output, err := runSplitCmd(t, src, "--out", out, "--sep", ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "wrote 2 file(s)") {
		t.Errorf("expected 2 files, got: %s", output)
	}
}

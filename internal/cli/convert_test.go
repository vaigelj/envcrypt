package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForConvert(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func runConvertCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envcrypt"}
	root.AddCommand(NewConvertCmd())
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"convert"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestConvertCmdToJSON(t *testing.T) {
	src := writeTempEnvForConvert(t, "FOO=bar\nBAZ=qux\n")
	out, err := runConvertCmd(t, src, "--from", "dotenv", "--to", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in JSON output, got: %s", out)
	}
}

func TestConvertCmdToShell(t *testing.T) {
	src := writeTempEnvForConvert(t, "FOO=bar\n")
	out, err := runConvertCmd(t, src, "--from", "dotenv", "--to", "shell")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export") {
		t.Errorf("expected shell export in output, got: %s", out)
	}
}

func TestConvertCmdToFile(t *testing.T) {
	src := writeTempEnvForConvert(t, "FOO=bar\n")
	dst := filepath.Join(t.TempDir(), "out.json")
	out, err := runConvertCmd(t, src, "--from", "dotenv", "--to", "json", "--output", dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "converted") {
		t.Errorf("expected confirmation message, got: %s", out)
	}
	data, _ := os.ReadFile(dst)
	if !strings.Contains(string(data), "FOO") {
		t.Errorf("expected FOO in output file, got: %s", data)
	}
}

func TestConvertCmdBadFormat(t *testing.T) {
	src := writeTempEnvForConvert(t, "FOO=bar\n")
	_, err := runConvertCmd(t, src, "--from", "dotenv", "--to", "yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

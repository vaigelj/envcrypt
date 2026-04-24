package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempEnvForFilter(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func runFilterCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	root := &cobra.Command{Use: "envcrypt"}
	root.AddCommand(NewFilterCmd())
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"filter"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestFilterCmdByPrefix(t *testing.T) {
	path := writeTempEnvForFilter(t, "DB_HOST=localhost\nAPP_NAME=myapp\nDB_PORT=5432\n")
	out, err := runFilterCmd("--prefix", "DB_", path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "DB_PORT") {
		t.Errorf("expected DB_ keys in output, got: %s", out)
	}
	if strings.Contains(out, "APP_NAME") {
		t.Errorf("APP_NAME should not appear in output")
	}
}

func TestFilterCmdExclude(t *testing.T) {
	path := writeTempEnvForFilter(t, "DB_HOST=localhost\nAPP_NAME=myapp\n")
	out, err := runFilterCmd("--prefix", "DB_", "--exclude", path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "DB_HOST") {
		t.Errorf("DB_HOST should be excluded")
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("APP_NAME should remain")
	}
}

func TestFilterCmdInPlace(t *testing.T) {
	path := writeTempEnvForFilter(t, "DB_HOST=localhost\nAPP_NAME=myapp\nDB_PORT=5432\n")
	_, err := runFilterCmd("--prefix", "APP_", "--in-place", path)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "DB_") {
		t.Errorf("DB_ keys should have been removed in-place")
	}
}

func TestFilterCmdPattern(t *testing.T) {
	path := writeTempEnvForFilter(t, "DB_HOST=localhost\nAPP_NAME=myapp\nSECRET_KEY=x\n")
	out, err := runFilterCmd("--pattern", `^(DB|SECRET)_`, path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected DB_HOST and SECRET_KEY, got: %s", out)
	}
	if strings.Contains(out, "APP_NAME") {
		t.Errorf("APP_NAME should not appear")
	}
}

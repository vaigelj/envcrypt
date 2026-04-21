package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

func writeTempEnvForBackup(t *testing.T, content string) string {
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

func runBackupCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := NewBackupCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestBackupCreateAndList(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("FOO=bar\nBAZ=qux\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Override working dir for backup storage
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	out, err := runBackupCmd(t, "create", envPath, "--label", "test-label")
	if err != nil {
		t.Fatalf("create: %v (out: %s)", err, out)
	}
	if !strings.Contains(out, "backup created:") {
		t.Errorf("unexpected output: %s", out)
	}

	out, err = runBackupCmd(t, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "test-label") {
		t.Errorf("expected label in list output, got: %s", out)
	}
}

func TestBackupRestoreAndDelete(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	entries := []envfile.Entry{{Key: "HELLO", Value: "world"}}
	b, err := envfile.CreateBackup(".", entries, "restore-test")
	if err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(dir, "restored.env")
	cmd := NewBackupCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"restore", b.ID, dest})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("restore: %v", err)
	}

	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "HELLO") {
		t.Errorf("restored file missing HELLO: %s", data)
	}

	delCmd := &cobra.Command{}
	_ = delCmd
	out, err := runBackupCmd(t, "delete", b.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if !strings.Contains(out, "deleted backup") {
		t.Errorf("unexpected delete output: %s", out)
	}
}

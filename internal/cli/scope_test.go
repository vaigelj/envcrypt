package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envcrypt/internal/envfile"
)

func runScopeCmd(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old) //nolint:errcheck

	cmd := NewScopeCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestScopeSetAndList(t *testing.T) {
	dir := t.TempDir()
	out, err := runScopeCmd(t, dir, "set", "web", "API_URL,PUBLIC_KEY")
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}
	if !strings.Contains(out, "web") {
		t.Errorf("expected scope name in output, got: %s", out)
	}

	out, err = runScopeCmd(t, dir, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "web") {
		t.Errorf("expected 'web' in list output, got: %s", out)
	}
}

func TestScopeShow(t *testing.T) {
	dir := t.TempDir()
	_ = envfile.AddScope(dir, "ops", []string{"LOG_LEVEL", "REGION"})

	out, err := runScopeCmd(t, dir, "show", "ops")
	if err != nil {
		t.Fatalf("show failed: %v", err)
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected LOG_LEVEL in output, got: %s", out)
	}
}

func TestScopeShowNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := runScopeCmd(t, dir, "show", "missing")
	if err == nil {
		t.Error("expected error for missing scope")
	}
}

func TestScopeRemove(t *testing.T) {
	dir := t.TempDir()
	_ = envfile.AddScope(dir, "tmp", []string{"X"})
	out, err := runScopeCmd(t, dir, "remove", "tmp")
	if err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected 'removed' in output, got: %s", out)
	}
	scopes, _ := envfile.LoadScopes(filepath.Join(dir))
	if len(scopes) != 0 {
		t.Errorf("expected 0 scopes after remove, got %d", len(scopes))
	}
}

func TestScopeListEmpty(t *testing.T) {
	dir := t.TempDir()
	out, err := runScopeCmd(t, dir, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "no scopes") {
		t.Errorf("expected 'no scopes' message, got: %s", out)
	}
}

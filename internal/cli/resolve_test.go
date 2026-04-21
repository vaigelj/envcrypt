package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempEnvForResolve(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func runResolveCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := NewResolveCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestResolveCmdBasic(t *testing.T) {
	path := writeTempEnvForResolve(t, "HOST=db\nURL=postgres://${HOST}/app\n")
	out, err := runResolveCmd(t, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "URL=postgres://db/app") {
		t.Errorf("expected resolved URL, got: %q", out)
	}
}

func TestResolveCmdStrictMissingFails(t *testing.T) {
	path := writeTempEnvForResolve(t, "URL=http://${GHOST}/path\n")
	_, err := runResolveCmd(t, "--strict", path)
	if err == nil {
		t.Fatal("expected error in strict mode")
	}
}

func TestResolveCmdLooseMissingOk(t *testing.T) {
	path := writeTempEnvForResolve(t, "URL=http://${GHOST}/path\n")
	out, err := runResolveCmd(t, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "${GHOST}") {
		t.Errorf("expected unresolved ref to remain, got: %q", out)
	}
}

func TestResolveCmdExportFormat(t *testing.T) {
	path := writeTempEnvForResolve(t, "NAME=world\n")
	out, err := runResolveCmd(t, "--format", "export", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "export NAME=") {
		t.Errorf("expected export prefix, got: %q", out)
	}
}

func TestResolveCmdUnknownFormat(t *testing.T) {
	path := writeTempEnvForResolve(t, "K=V\n")
	_, err := runResolveCmd(t, "--format", "yaml", path)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

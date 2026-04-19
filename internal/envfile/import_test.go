package envfile

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempImport(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "import*")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestImportDotenv(t *testing.T) {
	p := writeTempImport(t, "FOO=bar\nBAZ=qux\n")
	m, err := Import(p, "dotenv")
	if err != nil {
		t.Fatal(err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestImportJSON(t *testing.T) {
	data, _ := json.Marshal(map[string]string{"KEY": "val", "NUM": "42"})
	p := writeTempImport(t, string(data))
	m, err := Import(p, "json")
	if err != nil {
		t.Fatal(err)
	}
	if m["KEY"] != "val" {
		t.Fatalf("expected KEY=val, got %v", m["KEY"])
	}
}

func TestImportShell(t *testing.T) {
	p := writeTempImport(t, "export FOO='hello'\nexport BAR=\"world\"\n# comment\n")
	m, err := Import(p, "shell")
	if err != nil {
		t.Fatal(err)
	}
	if m["FOO"] != "hello" || m["BAR"] != "world" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestImportUnknownFormat(t *testing.T) {
	p := writeTempImport(t, "FOO=bar")
	_, err := Import(p, "yaml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestImportMissingFile(t *testing.T) {
	_, err := Import("/nonexistent/path.env", "dotenv")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

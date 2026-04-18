package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.template")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseTemplateBasic(t *testing.T) {
	p := writeTempTemplate(t, "DATABASE_URL= # db url\nAPI_KEY=\nOPTIONAL=\n")
	entries, err := ParseTemplate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "DATABASE_URL" || entries[0].Comment != "db url" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
}

func TestParseTemplateRequired(t *testing.T) {
	p := writeTempTemplate(t, "#! SECRET_KEY= # must be set\nOPTIONAL=\n")
	entries, err := ParseTemplate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !entries[0].Required {
		t.Error("expected SECRET_KEY to be required")
	}
	if entries[1].Required {
		t.Error("expected OPTIONAL to not be required")
	}
}

func TestParseTemplateInvalidLine(t *testing.T) {
	p := writeTempTemplate(t, "BADLINE\n")
	_, err := ParseTemplate(p)
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseTemplateMissingFile(t *testing.T) {
	_, err := ParseTemplate("/nonexistent/.env.template")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheckTemplate(t *testing.T) {
	entries := []TemplateEntry{
		{Key: "DB_URL", Required: true},
		{Key: "API_KEY", Required: true},
		{Key: "OPTIONAL", Required: false},
	}

	env := map[string]string{"DB_URL": "postgres://localhost", "OPTIONAL": "yes"}
	missing := CheckTemplate(entries, env)
	if len(missing) != 1 || missing[0] != "API_KEY" {
		t.Errorf("expected [API_KEY] missing, got %v", missing)
	}
}

func TestCheckTemplateAllPresent(t *testing.T) {
	entries := []TemplateEntry{
		{Key: "A", Required: true},
		{Key: "B", Required: true},
	}
	env := map[string]string{"A": "1", "B": "2"}
	missing := CheckTemplate(entries, env)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

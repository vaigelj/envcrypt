package envfile

import (
	"testing"
)

func TestResolvePlaceholdersBasic(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "DSN", Value: "postgres://{{HOST}}:5432/mydb"},
	}
	env := map[string]string{}
	got, err := ResolvePlaceholders(entries, env, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[1].Value != "postgres://localhost:5432/mydb" {
		t.Errorf("expected resolved DSN, got %q", got[1].Value)
	}
}

func TestResolvePlaceholdersExternalEnv(t *testing.T) {
	entries := []Entry{
		{Key: "GREETING", Value: "Hello, {{NAME}}!"},
	}
	env := map[string]string{"NAME": "World"}
	got, err := ResolvePlaceholders(entries, env, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0].Value != "Hello, World!" {
		t.Errorf("got %q", got[0].Value)
	}
}

func TestResolvePlaceholdersUnresolvedNonStrict(t *testing.T) {
	entries := []Entry{
		{Key: "MSG", Value: "Value is {{MISSING}}"},
	}
	got, err := ResolvePlaceholders(entries, map[string]string{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// placeholder should be left intact
	if got[0].Value != "Value is {{MISSING}}" {
		t.Errorf("got %q", got[0].Value)
	}
}

func TestResolvePlaceholdersUnresolvedStrict(t *testing.T) {
	entries := []Entry{
		{Key: "MSG", Value: "Value is {{MISSING}}"},
	}
	_, err := ResolvePlaceholders(entries, map[string]string{}, true)
	if err == nil {
		t.Fatal("expected error in strict mode for unresolved placeholder")
	}
}

func TestResolvePlaceholdersString(t *testing.T) {
	env := map[string]string{"APP": "envcrypt", "VER": "1.0"}
	result, err := ResolvePlaceholdersString("{{APP}} v{{VER}}", env, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "envcrypt v1.0" {
		t.Errorf("got %q", result)
	}
}

func TestResolvePlaceholdersMultiple(t *testing.T) {
	entries := []Entry{
		{Key: "PROTO", Value: "https"},
		{Key: "DOMAIN", Value: "example.com"},
		{Key: "URL", Value: "{{PROTO}}://{{DOMAIN}}/api"},
	}
	got, err := ResolvePlaceholders(entries, map[string]string{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[2].Value != "https://example.com/api" {
		t.Errorf("got %q", got[2].Value)
	}
}

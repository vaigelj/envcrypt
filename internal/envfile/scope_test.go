package envfile

import (
	"os"
	"testing"
)

func TestAddAndGetScope(t *testing.T) {
	dir := t.TempDir()
	if err := AddScope(dir, "frontend", []string{"API_URL", "PUBLIC_KEY"}); err != nil {
		t.Fatal(err)
	}
	scopes, err := LoadScopes(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 1 {
		t.Fatalf("expected 1 scope, got %d", len(scopes))
	}
	if scopes[0].Name != "frontend" {
		t.Errorf("unexpected name: %s", scopes[0].Name)
	}
	if len(scopes[0].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(scopes[0].Keys))
	}
}

func TestAddScopeOverwrites(t *testing.T) {
	dir := t.TempDir()
	_ = AddScope(dir, "backend", []string{"DB_URL"})
	_ = AddScope(dir, "backend", []string{"DB_URL", "SECRET"})
	scopes, _ := LoadScopes(dir)
	if len(scopes) != 1 {
		t.Fatalf("expected 1 scope after overwrite, got %d", len(scopes))
	}
	if len(scopes[0].Keys) != 2 {
		t.Errorf("expected 2 keys after overwrite, got %d", len(scopes[0].Keys))
	}
}

func TestRemoveScope(t *testing.T) {
	dir := t.TempDir()
	_ = AddScope(dir, "a", []string{"K1"})
	_ = AddScope(dir, "b", []string{"K2"})
	if err := RemoveScope(dir, "a"); err != nil {
		t.Fatal(err)
	}
	scopes, _ := LoadScopes(dir)
	if len(scopes) != 1 || scopes[0].Name != "b" {
		t.Errorf("unexpected scopes after remove: %+v", scopes)
	}
}

func TestLoadScopesMissing(t *testing.T) {
	dir := t.TempDir()
	scopes, err := LoadScopes(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 0 {
		t.Errorf("expected empty scopes, got %d", len(scopes))
	}
}

func TestApplyScope(t *testing.T) {
	dir := t.TempDir()
	_ = AddScope(dir, "web", []string{"API_URL", "PUBLIC_KEY"})
	entries := []Entry{
		{Key: "API_URL", Value: "http://example.com"},
		{Key: "SECRET", Value: "s3cr3t"},
		{Key: "PUBLIC_KEY", Value: "pk_abc"},
	}
	out, err := ApplyScope(dir, "web", entries)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 filtered entries, got %d", len(out))
	}
}

func TestApplyScopeNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := ApplyScope(dir, "nonexistent", nil)
	if err == nil {
		t.Error("expected error for missing scope")
	}
}

func TestScopePersistence(t *testing.T) {
	dir := t.TempDir()
	_ = AddScope(dir, "ops", []string{"LOG_LEVEL", "REGION"})
	// Re-load from disk
	scopes, err := LoadScopes(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 1 || scopes[0].Name != "ops" {
		t.Errorf("persistence failed: %+v", scopes)
	}
	_ = os.RemoveAll(dir)
}

package envfile

import (
	"os"
	"testing"
)

func TestValidateClean(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	verr, warnings := Validate(env)
	if verr != nil {
		t.Fatalf("expected no error, got: %v", verr)
	}
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got: %v", warnings)
	}
}

func TestValidateInvalidKey(t *testing.T) {
	env := map[string]string{
		"123BAD":  "value",
		"GOOD_KEY": "value",
	}
	verr, _ := Validate(env)
	if verr == nil {
		t.Fatal("expected validation error for invalid key")
	}
	if len(verr.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(verr.Issues), verr.Issues)
	}
}

func TestValidateEmptyValueWarning(t *testing.T) {
	env := map[string]string{
		"EMPTY_VAR": "",
		"SET_VAR":   "hello",
	}
	verr, warnings := Validate(env)
	if verr != nil {
		t.Fatalf("unexpected error: %v", verr)
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
}

func TestValidateFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("VALID_KEY=hello\nANOTHER=world\n")
	f.Close()

	verr, warnings, err := ValidateFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verr != nil {
		t.Fatalf("unexpected validation error: %v", verr)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
}

func TestValidateFileBadKey(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("GOOD=ok\n9BAD=nope\n")
	f.Close()

	verr, _, err := ValidateFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if verr == nil {
		t.Fatal("expected validation error")
	}
}

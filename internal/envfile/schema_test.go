package envfile

import (
	"testing"
)

func TestSchemaValidateClean(t *testing.T) {
	s := &Schema{
		Fields: []SchemaField{
			{Key: "DATABASE_URL", Required: true},
			{Key: "PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	env := map[string]string{"DATABASE_URL": "postgres://localhost/db", "PORT": "5432"}
	errs := s.Validate(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestSchemaValidateMissingRequired(t *testing.T) {
	s := &Schema{
		Fields: []SchemaField{
			{Key: "SECRET_KEY", Required: true},
		},
	}
	errs := s.Validate(map[string]string{})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Key != "SECRET_KEY" {
		t.Errorf("unexpected key: %s", errs[0].Key)
	}
}

func TestSchemaValidatePatternMismatch(t *testing.T) {
	s := &Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	errs := s.Validate(map[string]string{"PORT": "not-a-number"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestSchemaValidateOptionalMissing(t *testing.T) {
	s := &Schema{
		Fields: []SchemaField{
			{Key: "OPTIONAL_VAR", Required: false},
		},
	}
	errs := s.Validate(map[string]string{})
	if len(errs) != 0 {
		t.Fatalf("expected no errors for optional missing key, got %v", errs)
	}
}

func TestSchemaValidateInvalidPattern(t *testing.T) {
	s := &Schema{
		Fields: []SchemaField{
			{Key: "FOO", Required: true, Pattern: `[invalid`},
		},
	}
	errs := s.Validate(map[string]string{"FOO": "bar"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for invalid pattern, got %d", len(errs))
	}
}

package envfile

import (
	"testing"
)

func makeTypecheckEntries() []Entry {
	return []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "RATIO", Value: "3.14"},
		{Key: "DEBUG", Value: "true"},
		{Key: "SITE_URL", Value: "https://example.com"},
		{Key: "CONTACT", Value: "admin@example.com"},
		{Key: "CODE", Value: "ABC-123"},
	}
}

func TestTypecheckAllValid(t *testing.T) {
	rules := []TypeRule{
		{Key: "PORT", Type: "int"},
		{Key: "RATIO", Type: "float"},
		{Key: "DEBUG", Type: "bool"},
		{Key: "SITE_URL", Type: "url"},
		{Key: "CONTACT", Type: "email"},
		{Key: "CODE", Type: "regex", Pattern: `^[A-Z]+-\d+$`},
	}
	violations := TypeCheck(makeTypecheckEntries(), rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestTypecheckInvalidInt(t *testing.T) {
	entries := []Entry{{Key: "PORT", Value: "not-a-number"}}
	rules := []TypeRule{{Key: "PORT", Type: "int"}}
	vs := TypeCheck(entries, rules)
	if len(vs) != 1 || vs[0].Key != "PORT" {
		t.Fatalf("expected one violation for PORT, got %v", vs)
	}
}

func TestTypecheckInvalidBool(t *testing.T) {
	entries := []Entry{{Key: "DEBUG", Value: "yes"}}
	rules := []TypeRule{{Key: "DEBUG", Type: "bool"}}
	vs := TypeCheck(entries, rules)
	if len(vs) != 1 {
		t.Fatalf("expected one violation, got %v", vs)
	}
}

func TestTypecheckInvalidURL(t *testing.T) {
	entries := []Entry{{Key: "SITE_URL", Value: "ftp://bad"}}
	rules := []TypeRule{{Key: "SITE_URL", Type: "url"}}
	vs := TypeCheck(entries, rules)
	if len(vs) != 1 {
		t.Fatalf("expected one violation, got %v", vs)
	}
}

func TestTypecheckRequiredMissing(t *testing.T) {
	entries := []Entry{}
	rules := []TypeRule{{Key: "SECRET", Type: "int", Required: true}}
	vs := TypeCheck(entries, rules)
	if len(vs) != 1 || vs[0].Message != "required key is missing" {
		t.Fatalf("expected missing-key violation, got %v", vs)
	}
}

func TestTypecheckOptionalMissing(t *testing.T) {
	entries := []Entry{}
	rules := []TypeRule{{Key: "OPTIONAL", Type: "int", Required: false}}
	vs := TypeCheck(entries, rules)
	if len(vs) != 0 {
		t.Fatalf("expected no violations for optional missing key, got %v", vs)
	}
}

func TestTypecheckRegexMismatch(t *testing.T) {
	entries := []Entry{{Key: "CODE", Value: "abc-123"}}
	rules := []TypeRule{{Key: "CODE", Type: "regex", Pattern: `^[A-Z]+-\d+$`}}
	vs := TypeCheck(entries, rules)
	if len(vs) != 1 {
		t.Fatalf("expected one violation, got %v", vs)
	}
}

func TestTypecheckViolationError(t *testing.T) {
	v := TypeViolation{Key: "PORT", Value: "abc", Rule: "int", Message: "expected integer"}
	if v.Error() == "" {
		t.Fatal("expected non-empty error string")
	}
}

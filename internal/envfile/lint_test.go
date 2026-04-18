package envfile

import (
	"os"
	"testing"
)

func TestLintClean(t *testing.T) {
	entries := []Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "API_KEY", Value: "secret"},
	}
	issues := Lint(entries)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestLintDuplicateKey(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "FOO", Value: "baz"},
	}
	issues := Lint(entries)
	if !hasRule(issues, RuleDuplicateKey) {
		t.Fatal("expected duplicate-key issue")
	}
}

func TestLintKeyNotUppercase(t *testing.T) {
	entries := []Entry{
		{Key: "myKey", Value: "val"},
	}
	issues := Lint(entries)
	if !hasRule(issues, RuleKeyNotUppercase) {
		t.Fatal("expected key-not-uppercase issue")
	}
}

func TestLintLeadingUnderscore(t *testing.T) {
	entries := []Entry{
		{Key: "_INTERNAL", Value: "val"},
	}
	issues := Lint(entries)
	if !hasRule(issues, RuleLeadingUnderscore) {
		t.Fatal("expected leading-underscore issue")
	}
}

func TestLintWhitespaceValue(t *testing.T) {
	entries := []Entry{
		{Key: "EMPTY", Value: "   "},
	}
	issues := Lint(entries)
	if !hasRule(issues, RuleWhitespaceValue) {
		t.Fatal("expected whitespace-value issue")
	}
}

func TestLintFile(t *testing.T) {
	f, err := os.CreateTemp("", "lint*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("good_key=value\nGOOD=ok\n")
	f.Close()

	issues, err := LintFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	// good_key is not uppercase
	if !hasRule(issues, RuleKeyNotUppercase) {
		t.Fatal("expected key-not-uppercase from file lint")
	}
}

func hasRule(issues []LintIssue, rule LintRule) bool {
	for _, i := range issues {
		if i.Rule == rule {
			return true
		}
	}
	return false
}

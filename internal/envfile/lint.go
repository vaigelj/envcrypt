package envfile

import (
	"fmt"
	"strings"
)

// LintRule represents a single lint rule identifier.
type LintRule string

const (
	RuleDuplicateKey   LintRule = "duplicate-key"
	RuleKeyNotUppercase LintRule = "key-not-uppercase"
	RuleLeadingUnderscore LintRule = "leading-underscore"
	RuleWhitespaceValue  LintRule = "whitespace-value"
)

// LintIssue describes a single linting problem found in an env file.
type LintIssue struct {
	Line  int
	Key   string
	Rule  LintRule
	Msg   string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("line %d [%s] %s: %s", i.Line, i.Rule, i.Key, i.Msg)
}

// Lint checks a slice of Entry values for common style and correctness issues.
func Lint(entries []Entry) []LintIssue {
	var issues []LintIssue
	seen := make(map[string]int)

	for idx, e := range entries {
		line := idx + 1

		// duplicate key
		if prev, ok := seen[e.Key]; ok {
			issues = append(issues, LintIssue{
				Line: line,
				Key:  e.Key,
				Rule: RuleDuplicateKey,
				Msg:  fmt.Sprintf("duplicate of key first seen at entry %d", prev),
			})
		} else {
			seen[e.Key] = line
		}

		// key not uppercase
		if e.Key != strings.ToUpper(e.Key) {
			issues = append(issues, LintIssue{
				Line: line,
				Key:  e.Key,
				Rule: RuleKeyNotUppercase,
				Msg:  "key should be uppercase",
			})
		}

		// leading underscore
		if strings.HasPrefix(e.Key, "_") {
			issues = append(issues, LintIssue{
				Line: line,
				Key:  e.Key,
				Rule: RuleLeadingUnderscore,
				Msg:  "key starts with underscore",
			})
		}

		// whitespace-only value
		if len(e.Value) > 0 && strings.TrimSpace(e.Value) == "" {
			issues = append(issues, LintIssue{
				Line: line,
				Key:  e.Key,
				Rule: RuleWhitespaceValue,
				Msg:  "value contains only whitespace",
			})
		}
	}

	return issues
}

// LintFile parses the file at path and runs Lint on its entries.
func LintFile(path string) ([]LintIssue, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return nil, err
	}
	return Lint(entries), nil
}

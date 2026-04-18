package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// validKeyRe matches POSIX-style env var names.
var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError holds all issues found during validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("env validation failed:\n  %s", strings.Join(e.Issues, "\n  "))
}

// Validate checks a map of env vars for common issues:
//   - invalid key names
//   - empty values (reported as warnings via returned slice)
//   - duplicate keys are impossible in a map, so not checked here
func Validate(env map[string]string) (*ValidationError, []string) {
	var issues []string
	var warnings []string

	for k, v := range env {
		if !validKeyRe.MatchString(k) {
			issues = append(issues, fmt.Sprintf("invalid key name: %q", k))
		}
		if v == "" {
			warnings = append(warnings, fmt.Sprintf("key %q has an empty value", k))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}, warnings
	}
	return nil, warnings
}

// ValidateFile parses the file at path and runs Validate on its contents.
func ValidateFile(path string) (*ValidationError, []string, error) {
	env, err := Parse(path)
	if err != nil {
		return nil, nil, fmt.Errorf("parse: %w", err)
	}
	verr, warnings := Validate(env)
	return verr, warnings, nil
}

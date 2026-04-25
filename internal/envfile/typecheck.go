package envfile

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TypeRule defines an expected type constraint for an env key.
type TypeRule struct {
	Key      string
	Type     string // "int", "float", "bool", "url", "email", "regex"
	Pattern  string // used when Type == "regex"
	Required bool
}

// TypeViolation describes a single type-check failure.
type TypeViolation struct {
	Key     string
	Value   string
	Rule    string
	Message string
}

func (v TypeViolation) Error() string {
	return fmt.Sprintf("key %q value %q: %s", v.Key, v.Value, v.Message)
}

var (
	urlRE   = regexp.MustCompile(`^https?://[^\s]+$`)
	emailRE = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

// TypeCheck validates entries against a slice of TypeRules.
// It returns all violations found.
func TypeCheck(entries []Entry, rules []TypeRule) []TypeViolation {
	index := make(map[string]string, len(entries))
	for _, e := range entries {
		index[e.Key] = e.Value
	}

	var violations []TypeViolation
	for _, r := range rules {
		val, exists := index[r.Key]
		if !exists {
			if r.Required {
				violations = append(violations, TypeViolation{
					Key:     r.Key,
					Rule:    r.Type,
					Message: "required key is missing",
				})
			}
			continue
		}
		if msg := checkType(val, r); msg != "" {
			violations = append(violations, TypeViolation{
				Key:     r.Key,
				Value:   val,
				Rule:    r.Type,
				Message: msg,
			})
		}
	}
	return violations
}

func checkType(val string, r TypeRule) string {
	switch strings.ToLower(r.Type) {
	case "int":
		if _, err := strconv.ParseInt(val, 10, 64); err != nil {
			return "expected integer"
		}
	case "float":
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return "expected float"
		}
	case "bool":
		if _, err := strconv.ParseBool(val); err != nil {
			return "expected bool (true/false/1/0)"
		}
	case "url":
		if !urlRE.MatchString(val) {
			return "expected http/https URL"
		}
	case "email":
		if !emailRE.MatchString(val) {
			return "expected email address"
		}
	case "regex":
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return fmt.Sprintf("invalid rule pattern: %v", err)
		}
		if !re.MatchString(val) {
			return fmt.Sprintf("value does not match pattern %q", r.Pattern)
		}
	}
	return ""
}

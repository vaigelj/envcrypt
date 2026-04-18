package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField describes an expected env variable.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern for value validation
}

// Schema holds a set of expected fields for an env file.
type Schema struct {
	Fields []SchemaField
}

// SchemaError represents a single schema violation.
type SchemaError struct {
	Key     string
	Message string
}

func (e SchemaError) Error() string {
	return fmt.Sprintf("schema violation [%s]: %s", e.Key, e.Message)
}

// Validate checks env vars against the schema and returns all violations.
func (s *Schema) Validate(env map[string]string) []SchemaError {
	var errs []SchemaError
	for _, field := range s.Fields {
		val, ok := env[field.Key]
		if !ok || strings.TrimSpace(val) == "" {
			if field.Required {
				errs = append(errs, SchemaError{Key: field.Key, Message: "required key is missing or empty"})
			}
			continue
		}
		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				errs = append(errs, SchemaError{Key: field.Key, Message: fmt.Sprintf("invalid pattern: %v", err)})
				continue
			}
			if !re.MatchString(val) {
				errs = append(errs, SchemaError{Key: field.Key, Message: fmt.Sprintf("value does not match pattern %q", field.Pattern)})
			}
		}
	}
	return errs
}

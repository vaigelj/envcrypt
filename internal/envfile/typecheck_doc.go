// Package envfile provides the TypeCheck function for validating env entry
// values against declared type rules.
//
// Supported types:
//
//   - int   – value must be a valid 64-bit integer
//   - float – value must be a valid 64-bit float
//   - bool  – value must be parseable as a boolean (true/false/1/0)
//   - url   – value must be an http or https URL
//   - email – value must be a basic email address
//   - regex – value must match the rule's Pattern field
//
// Rules can be marked Required; a missing required key is reported as a
// TypeViolation even when no value is present.
package envfile

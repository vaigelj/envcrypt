// Package envfile provides utilities for parsing, writing, and manipulating
// .env files used by the envcrypt CLI.
//
// # Scopes
//
// A Scope is a named collection of environment variable keys that restricts
// which variables are visible or editable in a given context.
//
// Typical usage:
//
//	// Define a "frontend" scope
//	envfile.AddScope(dir, "frontend", []string{"API_URL", "PUBLIC_KEY"})
//
//	// Filter entries to only those in the scope
//	filtered, err := envfile.ApplyScope(dir, "frontend", allEntries)
//
// Scopes are persisted as JSON under .envcrypt/scopes.json.
package envfile

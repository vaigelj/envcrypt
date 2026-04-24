// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// # Flatten
//
// The Flatten family of functions converts structured or dot-notation keys
// into flat environment-variable-style keys.
//
// Example — flatten a slice of entries:
//
//	result := envfile.Flatten(entries, envfile.FlattenOptions{
//		Separator: "_",
//		Uppercase: true,
//		Prefix:    "APP_",
//	})
//
// Example — reconstruct a nested map from flat entries:
//
//	nested := envfile.UnflattenToMap(entries, "_")
package envfile

// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// # Filter
//
// The Filter function returns a subset of [Entry] values based on configurable
// criteria such as exact key match, key prefix, key suffix, or a regular
// expression pattern. Filtering is inclusive by default; pass
// [WithFilterExclude] to invert the selection and remove matched entries
// instead.
//
// Example — keep only keys that start with "DB_":
//
//	result, err := envfile.Filter(entries, envfile.WithFilterPrefix("DB_"))
//
// Example — remove all keys ending with "_SECRET":
//
//	result, err := envfile.Filter(entries,
//	    envfile.WithFilterSuffix("_SECRET"),
//	    envfile.WithFilterExclude(),
//	)
package envfile

// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// # CopyEnv / CopyFile
//
// CopyEnv merges key/value pairs from one map into another, with optional
// overwrite and exclusion controls. CopyFile is a convenience wrapper that
// reads source and destination files from disk, performs the copy, and writes
// the result back.
//
// Example:
//
//	n, err := envfile.CopyFile(".env", ".env.example", envfile.CopyOptions{
//		Overwrite: false,
//		Exclude:   []string{"SECRET_KEY"},
//	})
package envfile

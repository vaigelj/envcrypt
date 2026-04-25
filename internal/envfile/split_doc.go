// Package envfile provides utilities for reading, writing, and manipulating
// .env files.
//
// # Split
//
// The Split and SplitFile functions partition environment variable entries into
// groups based on their key prefix. The prefix is determined by splitting each
// key on a configurable separator (default "_") and taking the first segment.
//
// Example:
//
//	DB_HOST=localhost   → group "DB"
//	DB_PORT=5432        → group "DB"
//	APP_NAME=envcrypt   → group "APP"
//	STANDALONE=yes      → group "_default"
//
// SplitFile writes one <prefix>.env file per group into a target directory.
package envfile

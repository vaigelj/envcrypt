// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// # Reorder
//
// Reorder rearranges the entries in a []Entry slice so that a caller-supplied
// list of keys appears first, in the specified order. Any entries whose keys
// are not mentioned in the order list follow in their original relative order.
//
// Example:
//
//	entries, _ := envfile.ParseFile(".env")
//	ordered, err := envfile.Reorder(entries, []string{"APP_ENV", "APP_PORT"})
//
// ReorderFile is a convenience wrapper that reads, reorders, and atomically
// writes a file in one call.
package envfile

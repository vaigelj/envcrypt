// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// # Resolve
//
// The Resolve function expands shell-style variable references (${VAR} or
// $VAR) found in entry values, using other entries in the same file as the
// source of truth.
//
// Two resolution modes are supported:
//
//   - ResolveModeStrict (default): returns an error when a referenced
//     variable cannot be found.
//   - ResolveModeLoose: leaves unresolvable references unchanged so the
//     caller can decide what to do.
//
// When ResolveOptions.Environ is true, OS environment variables are
// consulted as a fallback after the entries map is checked.
package envfile

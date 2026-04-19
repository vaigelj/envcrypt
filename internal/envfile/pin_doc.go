// Package envfile provides utilities for parsing, writing, and managing
// .env files.
//
// Pin support allows teams to capture immutable named snapshots of an env
// file at a point in time. Pins are stored as JSON files under a configurable
// directory and can be listed, retrieved, or deleted via the CLI.
//
// Example usage:
//
//	envfile.SavePin("/home/user/.envcrypt", "release-1.2", values)
//	pin, _ := envfile.LoadPin("/home/user/.envcrypt", "release-1.2")
//	names, _ := envfile.ListPins("/home/user/.envcrypt")
package envfile

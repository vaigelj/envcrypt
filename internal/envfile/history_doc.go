// Package envfile provides utilities for parsing, writing, and managing
// .env files.
//
// # History
//
// The history sub-feature records versioned snapshots of an env file on disk
// as a JSON sidecar (<file>.history.json). Each entry captures:
//
//   - Timestamp – when the snapshot was taken (UTC)
//   - Label     – a human-readable tag (e.g. "pre-rotation", "v2")
//   - Values    – the full key/value map at that point in time
//
// Use AppendHistory to record a new entry, LoadHistory to read all entries,
// and ClearHistory to remove the sidecar file entirely.
package envfile

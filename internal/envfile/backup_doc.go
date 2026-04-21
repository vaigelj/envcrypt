// Package envfile provides utilities for managing .env files.
//
// # Backup
//
// The backup sub-feature allows callers to snapshot the current state of an
// env file's entries and restore them later.  Backups are stored as JSON
// files under <dir>/.envcrypt/backups/ and identified by a nanosecond
// timestamp ID.
//
// Usage:
//
//	b, err := envfile.CreateBackup(dir, entries, "before-rotation")
//	list, err := envfile.ListBackups(dir)
//	loaded, err := envfile.LoadBackup(dir, b.ID)
//	err = envfile.DeleteBackup(dir, b.ID)
package envfile

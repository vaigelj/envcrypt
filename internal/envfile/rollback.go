// Package envfile provides utilities for parsing, writing, and manipulating .env files.
package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// RollbackEntry represents a single rollback checkpoint for an env file.
type RollbackEntry struct {
	ID        string            `json:"id"`
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Entries   []Entry           `json:"entries"`
	Meta      map[string]string `json:"meta,omitempty"`
}

func rollbackDir(dir string) string {
	return filepath.Join(dir, ".envcrypt", "rollbacks")
}

func rollbackPath(dir, id string) string {
	return filepath.Join(rollbackDir(dir), id+".json")
}

// SaveRollback persists the given entries as a named rollback checkpoint.
// dir is the directory where the .envcrypt folder resides.
// label is a human-readable description (may be empty).
func SaveRollback(dir, label string, entries []Entry) (RollbackEntry, error) {
	if err := os.MkdirAll(rollbackDir(dir), 0700); err != nil {
		return RollbackEntry{}, fmt.Errorf("rollback: mkdir: %w", err)
	}

	now := time.Now().UTC()
	id := fmt.Sprintf("%d", now.UnixNano())

	re := RollbackEntry{
		ID:        id,
		Label:     label,
		CreatedAt: now,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(re, "", "  ")
	if err != nil {
		return RollbackEntry{}, fmt.Errorf("rollback: marshal: %w", err)
	}

	if err := os.WriteFile(rollbackPath(dir, id), data, 0600); err != nil {
		return RollbackEntry{}, fmt.Errorf("rollback: write: %w", err)
	}

	return re, nil
}

// LoadRollback loads a rollback checkpoint by ID.
func LoadRollback(dir, id string) (RollbackEntry, error) {
	data, err := os.ReadFile(rollbackPath(dir, id))
	if err != nil {
		if os.IsNotExist(err) {
			return RollbackEntry{}, fmt.Errorf("rollback %q not found", id)
		}
		return RollbackEntry{}, fmt.Errorf("rollback: read: %w", err)
	}

	var re RollbackEntry
	if err := json.Unmarshal(data, &re); err != nil {
		return RollbackEntry{}, fmt.Errorf("rollback: unmarshal: %w", err)
	}

	return re, nil
}

// ListRollbacks returns all rollback checkpoints in the given directory,
// sorted from newest to oldest.
func ListRollbacks(dir string) ([]RollbackEntry, error) {
	entries, err := os.ReadDir(rollbackDir(dir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("rollback: readdir: %w", err)
	}

	var result []RollbackEntry
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		re, err := LoadRollback(dir, id)
		if err != nil {
			continue
		}
		result = append(result, re)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

// DeleteRollback removes a rollback checkpoint by ID.
func DeleteRollback(dir, id string) error {
	path := rollbackPath(dir, id)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("rollback %q not found", id)
		}
		return fmt.Errorf("rollback: delete: %w", err)
	}
	return nil
}

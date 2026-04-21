package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// BackupEntry represents a single backup of an env file.
type BackupEntry struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label,omitempty"`
	Entries   []Entry           `json:"entries"`
	Meta      map[string]string `json:"meta,omitempty"`
}

func backupDir(dir string) string {
	return filepath.Join(dir, ".envcrypt", "backups")
}

func backupPath(dir, id string) string {
	return filepath.Join(backupDir(dir), id+".json")
}

// CreateBackup saves the current entries as a named backup under dir.
func CreateBackup(dir string, entries []Entry, label string) (BackupEntry, error) {
	if err := os.MkdirAll(backupDir(dir), 0700); err != nil {
		return BackupEntry{}, fmt.Errorf("backup: mkdir: %w", err)
	}
	b := BackupEntry{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return BackupEntry{}, fmt.Errorf("backup: marshal: %w", err)
	}
	if err := os.WriteFile(backupPath(dir, b.ID), data, 0600); err != nil {
		return BackupEntry{}, fmt.Errorf("backup: write: %w", err)
	}
	return b, nil
}

// ListBackups returns all backups stored under dir, sorted newest first.
func ListBackups(dir string) ([]BackupEntry, error) {
	glob := filepath.Join(backupDir(dir), "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("backup: glob: %w", err)
	}
	var result []BackupEntry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			continue
		}
		var b BackupEntry
		if err := json.Unmarshal(data, &b); err != nil {
			continue
		}
		result = append(result, b)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	return result, nil
}

// LoadBackup retrieves a single backup by ID from dir.
func LoadBackup(dir, id string) (BackupEntry, error) {
	data, err := os.ReadFile(backupPath(dir, id))
	if err != nil {
		if os.IsNotExist(err) {
			return BackupEntry{}, fmt.Errorf("backup %q not found", id)
		}
		return BackupEntry{}, fmt.Errorf("backup: read: %w", err)
	}
	var b BackupEntry
	if err := json.Unmarshal(data, &b); err != nil {
		return BackupEntry{}, fmt.Errorf("backup: unmarshal: %w", err)
	}
	return b, nil
}

// DeleteBackup removes a backup by ID from dir.
func DeleteBackup(dir, id string) error {
	err := os.Remove(backupPath(dir, id))
	if os.IsNotExist(err) {
		return fmt.Errorf("backup %q not found", id)
	}
	return err
}

package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ArchiveEntry holds a named, timestamped snapshot of env entries.
type ArchiveEntry struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Entries   []Entry   `json:"entries"`
}

func archiveDir(base string) string {
	return filepath.Join(base, ".envcrypt", "archives")
}

func archivePath(base, name string) string {
	return filepath.Join(archiveDir(base), name+".json")
}

// SaveArchive persists entries under the given name inside base directory.
func SaveArchive(base, name string, entries []Entry) error {
	if name == "" {
		return fmt.Errorf("archive name must not be empty")
	}
	dir := archiveDir(base)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create archive dir: %w", err)
	}
	a := ArchiveEntry{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal archive: %w", err)
	}
	return os.WriteFile(archivePath(base, name), data, 0o600)
}

// LoadArchive retrieves a named archive from base directory.
func LoadArchive(base, name string) (ArchiveEntry, error) {
	data, err := os.ReadFile(archivePath(base, name))
	if err != nil {
		if os.IsNotExist(err) {
			return ArchiveEntry{}, fmt.Errorf("archive %q not found", name)
		}
		return ArchiveEntry{}, fmt.Errorf("read archive: %w", err)
	}
	var a ArchiveEntry
	if err := json.Unmarshal(data, &a); err != nil {
		return ArchiveEntry{}, fmt.Errorf("unmarshal archive: %w", err)
	}
	return a, nil
}

// ListArchives returns all archive names sorted alphabetically.
func ListArchives(base string) ([]string, error) {
	dir := archiveDir(base)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list archives: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	sort.Strings(names)
	return names, nil
}

// DeleteArchive removes a named archive.
func DeleteArchive(base, name string) error {
	err := os.Remove(archivePath(base, name))
	if os.IsNotExist(err) {
		return fmt.Errorf("archive %q not found", name)
	}
	return err
}

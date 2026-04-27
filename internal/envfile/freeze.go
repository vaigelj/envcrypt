package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FrozenEnv represents a read-only snapshot of an env file with metadata.
type FrozenEnv struct {
	FrozenAt time.Time        `json:"frozen_at"`
	Source   string           `json:"source"`
	Entries  []Entry          `json:"entries"`
	Checksum string           `json:"checksum"`
}

func freezePath(dir, name string) string {
	return filepath.Join(dir, ".envcrypt", "frozen", name+".frozen.json")
}

// Freeze captures the current state of entries as immutable and saves it.
func Freeze(dir, name, source string, entries []Entry) (*FrozenEnv, error) {
	f := &FrozenEnv{
		FrozenAt: time.Now().UTC(),
		Source:   source,
		Entries:  entries,
		Checksum: entriesChecksum(entries),
	}
	path := freezePath(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("freeze: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("freeze: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return nil, fmt.Errorf("freeze: write: %w", err)
	}
	return f, nil
}

// LoadFrozen reads a previously frozen env by name.
func LoadFrozen(dir, name string) (*FrozenEnv, error) {
	path := freezePath(dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("freeze: %q not found", name)
		}
		return nil, fmt.Errorf("freeze: read: %w", err)
	}
	var f FrozenEnv
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("freeze: unmarshal: %w", err)
	}
	return &f, nil
}

// IsTampered returns true if the entries no longer match the stored checksum.
func (f *FrozenEnv) IsTampered() bool {
	return entriesChecksum(f.Entries) != f.Checksum
}

// DeleteFrozen removes a frozen env file.
func DeleteFrozen(dir, name string) error {
	path := freezePath(dir, name)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("freeze: delete: %w", err)
	}
	return nil
}

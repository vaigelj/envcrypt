package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Entries   map[string]string `json:"entries"`
}

// TakeSnapshot reads the given env file and returns a Snapshot.
func TakeSnapshot(path string) (*Snapshot, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: %w", err)
	}
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    path,
		Entries:   m,
	}, nil
}

// SaveSnapshot writes a Snapshot to a JSON file at dest.
func SaveSnapshot(snap *Snapshot, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("snapshot save: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// LoadSnapshot reads a Snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot load: %w", err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot decode: %w", err)
	}
	return &snap, nil
}

// DiffSnapshot compares two snapshots and returns a ChangeSet.
func DiffSnapshot(old, new *Snapshot) []Change {
	oldEntries := make([]Entry, 0, len(old.Entries))
	for k, v := range old.Entries {
		oldEntries = append(oldEntries, Entry{Key: k, Value: v})
	}
	newEntries := make([]Entry, 0, len(new.Entries))
	for k, v := range new.Entries {
		newEntries = append(newEntries, Entry{Key: k, Value: v})
	}
	return Compare(oldEntries, newEntries)
}

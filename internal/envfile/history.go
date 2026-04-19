package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a single versioned snapshot of an env file.
type HistoryEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label"`
	Values    map[string]string `json:"values"`
}

// HistoryFile is the on-disk representation of all history entries.
type HistoryFile struct {
	Entries []HistoryEntry `json:"entries"`
}

// historyPath returns the path to the history file for a given env file.
func historyPath(envPath string) string {
	return envPath + ".history.json"
}

// AppendHistory adds a new entry to the history for the given env file.
func AppendHistory(envPath, label string, values map[string]string) error {
	hf, err := LoadHistory(envPath)
	if err != nil {
		hf = &HistoryFile{}
	}
	hf.Entries = append(hf.Entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Values:    values,
	})
	data, err := json.MarshalIndent(hf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(historyPath(envPath), data, 0600)
}

// LoadHistory reads all history entries for the given env file.
func LoadHistory(envPath string) (*HistoryFile, error) {
	data, err := os.ReadFile(historyPath(envPath))
	if err != nil {
		return nil, err
	}
	var hf HistoryFile
	if err := json.Unmarshal(data, &hf); err != nil {
		return nil, err
	}
	return &hf, nil
}

// ClearHistory removes the history file for the given env file.
func ClearHistory(envPath string) error {
	p := historyPath(envPath)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(p)
}

// HistoryDir returns all history files found in a directory.
func HistoryDir(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.history.json"))
	if err != nil {
		return nil, err
	}
	return matches, nil
}

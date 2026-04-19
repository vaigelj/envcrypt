package envfile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Pin records a named, immutable snapshot of an env file.
type Pin struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Values    map[string]string `json:"values"`
}

func pinPath(dir, name string) string {
	return filepath.Join(dir, "pins", name+".json")
}

// SavePin persists a named pin under dir.
func SavePin(dir, name string, values map[string]string) error {
	p := pinPath(dir, name)
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	pin := Pin{Name: name, CreatedAt: time.Now().UTC(), Values: values}
	data, err := json.MarshalIndent(pin, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

// LoadPin retrieves a named pin from dir.
func LoadPin(dir, name string) (*Pin, error) {
	data, err := os.ReadFile(pinPath(dir, name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var pin Pin
	if err := json.Unmarshal(data, &pin); err != nil {
		return nil, err
	}
	return &pin, nil
}

// ListPins returns all pin names stored under dir.
func ListPins(dir string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(dir, "pins"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			n := e.Name()
			names = append(names, n[:len(n)-5]) // strip .json
		}
	}
	return names, nil
}

// DeletePin removes a named pin.
func DeletePin(dir, name string) error {
	return os.Remove(pinPath(dir, name))
}

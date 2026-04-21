package envfile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Scope represents a named set of key filters that restrict which env keys
// are visible or editable in a given context (e.g. "frontend", "backend").
type Scope struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

func scopesPath(dir string) string {
	return filepath.Join(dir, ".envcrypt", "scopes.json")
}

// LoadScopes reads all scopes from the given directory.
func LoadScopes(dir string) ([]Scope, error) {
	path := scopesPath(dir)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Scope{}, nil
		}
		return nil, err
	}
	var scopes []Scope
	if err := json.Unmarshal(data, &scopes); err != nil {
		return nil, err
	}
	return scopes, nil
}

// SaveScopes writes the given scopes to disk.
func SaveScopes(dir string, scopes []Scope) error {
	path := scopesPath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(scopes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// AddScope adds or replaces a scope by name.
func AddScope(dir, name string, keys []string) error {
	scopes, err := LoadScopes(dir)
	if err != nil {
		return err
	}
	for i, s := range scopes {
		if s.Name == name {
			scopes[i].Keys = keys
			return SaveScopes(dir, scopes)
		}
	}
	scopes = append(scopes, Scope{Name: name, Keys: keys})
	return SaveScopes(dir, scopes)
}

// RemoveScope deletes a scope by name.
func RemoveScope(dir, name string) error {
	scopes, err := LoadScopes(dir)
	if err != nil {
		return err
	}
	filtered := scopes[:0]
	for _, s := range scopes {
		if s.Name != name {
			filtered = append(filtered, s)
		}
	}
	return SaveScopes(dir, filtered)
}

// ApplyScope filters entries to only those whose key appears in the scope.
// If the scope name is not found, an error is returned.
func ApplyScope(dir, name string, entries []Entry) ([]Entry, error) {
	scopes, err := LoadScopes(dir)
	if err != nil {
		return nil, err
	}
	for _, s := range scopes {
		if s.Name == name {
			allowed := make(map[string]bool, len(s.Keys))
			for _, k := range s.Keys {
				allowed[k] = true
			}
			var out []Entry
			for _, e := range entries {
				if allowed[e.Key] {
					out = append(out, e)
				}
			}
			return out, nil
		}
	}
	return nil, errors.New("scope not found: " + name)
}

package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Tag represents a named label attached to a set of env keys.
type Tag struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// TagStore holds all tags for an env file.
type TagStore struct {
	Tags []Tag `json:"tags"`
}

func tagsPath(dir string) string {
	return filepath.Join(dir, ".envcrypt_tags.json")
}

// LoadTags loads the tag store from dir. Returns empty store if missing.
func LoadTags(dir string) (*TagStore, error) {
	data, err := os.ReadFile(tagsPath(dir))
	if os.IsNotExist(err) {
		return &TagStore{}, nil
	}
	if err != nil {
		return nil, err
	}
	var ts TagStore
	if err := json.Unmarshal(data, &ts); err != nil {
		return nil, err
	}
	return &ts, nil
}

// SaveTags persists the tag store to dir.
func SaveTags(dir string, ts *TagStore) error {
	data, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tagsPath(dir), data, 0600)
}

// AddTag adds or replaces a tag by name.
func (ts *TagStore) AddTag(name string, keys []string) {
	for i, t := range ts.Tags {
		if t.Name == name {
			ts.Tags[i].Keys = keys
			return
		}
	}
	ts.Tags = append(ts.Tags, Tag{Name: name, Keys: keys})
}

// RemoveTag deletes a tag by name. Returns false if not found.
func (ts *TagStore) RemoveTag(name string) bool {
	for i, t := range ts.Tags {
		if t.Name == name {
			ts.Tags = append(ts.Tags[:i], ts.Tags[i+1:]...)
			return true
		}
	}
	return false
}

// GetTag returns the tag with the given name, or nil.
func (ts *TagStore) GetTag(name string) *Tag {
	for _, t := range ts.Tags {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

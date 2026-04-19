package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Group represents a named collection of env var keys.
type Group struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

func groupsPath(dir string) string {
	return filepath.Join(dir, ".envcrypt_groups.json")
}

// LoadGroups loads all groups from the given directory.
func LoadGroups(dir string) ([]Group, error) {
	path := groupsPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Group{}, nil
	}
	if err != nil {
		return nil, err
	}
	var groups []Group
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// SaveGroups persists groups to the given directory.
func SaveGroups(dir string, groups []Group) error {
	data, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(groupsPath(dir), data, 0600)
}

// AddGroup adds or replaces a group by name.
func AddGroup(dir, name string, keys []string) error {
	groups, err := LoadGroups(dir)
	if err != nil {
		return err
	}
	for i, g := range groups {
		if g.Name == name {
			groups[i].Keys = keys
			return SaveGroups(dir, groups)
		}
	}
	groups = append(groups, Group{Name: name, Keys: keys})
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return SaveGroups(dir, groups)
}

// RemoveGroup deletes a group by name.
func RemoveGroup(dir, name string) error {
	groups, err := LoadGroups(dir)
	if err != nil {
		return err
	}
	filtered := groups[:0]
	for _, g := range groups {
		if g.Name != name {
			filtered = append(filtered, g)
		}
	}
	if len(filtered) == len(groups) {
		return fmt.Errorf("group %q not found", name)
	}
	return SaveGroups(dir, filtered)
}

// GetGroup returns a group by name.
func GetGroup(dir, name string) (Group, error) {
	groups, err := LoadGroups(dir)
	if err != nil {
		return Group{}, err
	}
	for _, g := range groups {
		if g.Name == name {
			return g, nil
		}
	}
	return Group{}, fmt.Errorf("group %q not found", name)
}

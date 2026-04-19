package envfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// Profile represents a named environment profile (e.g. "dev", "staging", "prod").
type Profile struct {
	Name string
	Path string
}

// ListProfiles scans dir for files matching .env.<profile> and returns profiles found.
func ListProfiles(dir string) ([]Profile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	var profiles []Profile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) > 5 && name[:5] == ".env." {
			profiles = append(profiles, Profile{
				Name: name[5:],
				Path: filepath.Join(dir, name),
			})
		}
	}
	return profiles, nil
}

// LoadProfile parses the env file for the given profile name from dir.
func LoadProfile(dir, profile string) (map[string]string, error) {
	path := filepath.Join(dir, ".env."+profile)
	pairs, err := ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("load profile %q: %w", profile, err)
	}
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		m[p.Key] = p.Value
	}
	return m, nil
}

// SaveProfile writes key/value pairs to .env.<profile> in dir.
func SaveProfile(dir, profile string, data map[string]string) error {
	path := filepath.Join(dir, ".env."+profile)
	pairs := make([]Pair, 0, len(data))
	for k, v := range data {
		pairs = append(pairs, Pair{Key: k, Value: v})
	}
	if err := WriteFile(path, pairs); err != nil {
		return fmt.Errorf("save profile %q: %w", profile, err)
	}
	return nil
}

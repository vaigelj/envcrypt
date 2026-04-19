package envfile

import "fmt"

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Renamed bool
}

// RenameKey renames a key in the given env map, returning a result.
// Returns an error if oldKey does not exist or newKey already exists.
func RenameKey(env map[string]string, oldKey, newKey string) (RenameResult, error) {
	if _, ok := env[oldKey]; !ok {
		return RenameResult{}, fmt.Errorf("key %q not found", oldKey)
	}
	if _, ok := env[newKey]; ok {
		return RenameResult{}, fmt.Errorf("key %q already exists", newKey)
	}
	env[newKey] = env[oldKey]
	delete(env, oldKey)
	return RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: true}, nil
}

// RenameFile renames a key in an env file on disk, preserving all other entries.
func RenameFile(path, oldKey, newKey string) (RenameResult, error) {
	env, err := ParseFile(path)
	if err != nil {
		return RenameResult{}, fmt.Errorf("parse: %w", err)
	}
	res, err := RenameKey(env, oldKey, newKey)
	if err != nil {
		return RenameResult{}, err
	}
	if err := WriteFile(path, env); err != nil {
		return RenameResult{}, fmt.Errorf("write: %w", err)
	}
	return res, nil
}

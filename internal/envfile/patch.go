package envfile

import "fmt"

// PatchOp represents a single patch operation kind.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
	PatchRename PatchOp = "rename"
)

// PatchInstruction describes one atomic change to apply to an env map.
type PatchInstruction struct {
	Op      PatchOp
	Key     string
	Value   string // used by PatchSet
	NewKey  string // used by PatchRename
}

// Patch applies a slice of PatchInstructions to entries and returns the
// modified slice. Operations are applied in order; an error stops processing.
func Patch(entries []Entry, instructions []PatchInstruction) ([]Entry, error) {
	result := make([]Entry, len(entries))
	copy(result, entries)

	for _, inst := range instructions {
		switch inst.Op {
		case PatchSet:
			result = applySet(result, inst.Key, inst.Value)
		case PatchDelete:
			result = applyDelete(result, inst.Key)
		case PatchRename:
			var err error
			result, err = applyRename(result, inst.Key, inst.NewKey)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("patch: unknown op %q", inst.Op)
		}
	}
	return result, nil
}

// PatchFile reads path, applies instructions, and writes the result back.
func PatchFile(path string, instructions []PatchInstruction) error {
	entries, err := ParseFile(path)
	if err != nil {
		return err
	}
	patched, err := Patch(entries, instructions)
	if err != nil {
		return err
	}
	return WriteFile(path, patched)
}

func applySet(entries []Entry, key, value string) []Entry {
	for i, e := range entries {
		if e.Key == key {
			entries[i].Value = value
			return entries
		}
	}
	return append(entries, Entry{Key: key, Value: value})
}

func applyDelete(entries []Entry, key string) []Entry {
	out := entries[:0:len(entries)]
	for _, e := range entries {
		if e.Key != key {
			out = append(out, e)
		}
	}
	return out
}

func applyRename(entries []Entry, oldKey, newKey string) ([]Entry, error) {
	found := false
	for i, e := range entries {
		if e.Key == oldKey {
			entries[i].Key = newKey
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("patch rename: key %q not found", oldKey)
	}
	return entries, nil
}

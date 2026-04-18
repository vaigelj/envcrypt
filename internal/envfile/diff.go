package envfile

// ChangeType describes what happened to a key between two env snapshots.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// Change represents a single key-level difference between two env maps.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// Compare returns the list of changes needed to go from base to updated.
func Compare(base, updated map[string]string) []Change {
	var changes []Change

	for k, newVal := range updated {
		oldVal, exists := base[k]
		if !exists {
			changes = append(changes, Change{Key: k, Type: ChangeAdded, NewVal: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Type: ChangeUpdated, OldVal: oldVal, NewVal: newVal})
		}
	}

	for k, oldVal := range base {
		if _, exists := updated[k]; !exists {
			changes = append(changes, Change{Key: k, Type: ChangeRemoved, OldVal: oldVal})
		}
	}

	return changes
}

// Summary returns counts of each change type.
func Summary(changes []Change) (added, updated, removed int) {
	for _, c := range changes {
		switch c.Type {
		case ChangeAdded:
			added++
		case ChangeUpdated:
			updated++
		case ChangeRemoved:
			removed++
		}
	}
	return
}

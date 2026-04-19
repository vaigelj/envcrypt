package envfile

import "fmt"

// Version represents a named snapshot of an env file.
type Version struct {
	Name  string
	Vars  map[string]string
}

// VersionDiff describes changes between two named versions.
type VersionDiff struct {
	From    string
	To      string
	Changes []DiffEntry
}

// DiffEntry describes a single key change between versions.
type DiffEntry struct {
	Key    string
	OldVal string
	NewVal string
	Op     string // added, removed, updated
}

// CompareVersions diffs two Version snapshots and returns a VersionDiff.
func CompareVersions(from, to Version) VersionDiff {
	var entries []DiffEntry
	for k, v := range from.Vars {
		if nv, ok := to.Vars[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, OldVal: v, Op: "removed"})
		} else if nv != v {
			entries = append(entries, DiffEntry{Key: k, OldVal: v, NewVal: nv, Op: "updated"})
		}
	}
	for k, v := range to.Vars {
		if _, ok := from.Vars[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, NewVal: v, Op: "added"})
		}
	}
	return VersionDiff{From: from.Name, To: to.Name, Changes: entries}
}

// FormatVersionDiff returns a human-readable summary of a VersionDiff.
func FormatVersionDiff(d VersionDiff) string {
	if len(d.Changes) == 0 {
		return fmt.Sprintf("No changes between %s and %s.", d.From, d.To)
	}
	out := fmt.Sprintf("Changes from %s to %s:\n", d.From, d.To)
	for _, e := range d.Changes {
		switch e.Op {
		case "added":
			out += fmt.Sprintf("  + %s=%s\n", e.Key, e.NewVal)
		case "removed":
			out += fmt.Sprintf("  - %s=%s\n", e.Key, e.OldVal)
		case "updated":
			out += fmt.Sprintf("  ~ %s: %s -> %s\n", e.Key, e.OldVal, e.NewVal)
		}
	}
	return out
}

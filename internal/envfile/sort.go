package envfile

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for env entries.
type SortOrder int

const (
	SortAlpha    SortOrder = iota // alphabetical by key
	SortAlphaDesc                 // reverse alphabetical
	SortByLength                  // shortest key first
)

// SortOptions configures sorting behaviour.
type SortOptions struct {
	Order   SortOrder
	Groups  []string // keep these key prefixes together at the top
}

// Sort returns a new slice of entries sorted according to opts.
func Sort(entries []Entry, opts SortOptions) []Entry {
	result := make([]Entry, len(entries))
	copy(result, entries)

	pinned := map[string]int{}
	for i, prefix := range opts.Groups {
		pinned[strings.ToUpper(prefix)] = i
	}

	groupOf := func(key string) int {
		for prefix, idx := range pinned {
			if strings.HasPrefix(strings.ToUpper(key), prefix) {
				return idx
			}
		}
		return len(opts.Groups)
	}

	sort.SliceStable(result, func(i, j int) bool {
		gi, gj := groupOf(result[i].Key), groupOf(result[j].Key)
		if gi != gj {
			return gi < gj
		}
		switch opts.Order {
		case SortAlphaDesc:
			return result[i].Key > result[j].Key
		case SortByLength:
			if len(result[i].Key) != len(result[j].Key) {
				return len(result[i].Key) < len(result[j].Key)
			}
			return result[i].Key < result[j].Key
		default: // SortAlpha
			return result[i].Key < result[j].Key
		}
	})
	return result
}

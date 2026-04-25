package envfile

import (
	"fmt"
	"os"
)

// ReorderOption configures how Reorder behaves.
type ReorderOption func(*reorderConfig)

type reorderConfig struct {
	missingOk bool
}

// WithReorderMissingOk allows keys in the order list that are absent from
// the entries to be silently skipped instead of returning an error.
func WithReorderMissingOk() ReorderOption {
	return func(c *reorderConfig) { c.missingOk = true }
}

// Reorder returns a new slice of Entry values arranged so that the keys
// listed in order appear first (in that order), followed by any remaining
// entries in their original relative order.
//
// If a key in order does not exist in entries and WithReorderMissingOk is
// not set, an error is returned.
func Reorder(entries []Entry, order []string, opts ...ReorderOption) ([]Entry, error) {
	cfg := &reorderConfig{}
	for _, o := range opts {
		o(cfg)
	}

	index := make(map[string]Entry, len(entries))
	for _, e := range entries {
		index[e.Key] = e
	}

	seen := make(map[string]bool, len(order))
	result := make([]Entry, 0, len(entries))

	for _, k := range order {
		e, ok := index[k]
		if !ok {
			if !cfg.missingOk {
				return nil, fmt.Errorf("reorder: key %q not found in entries", k)
			}
			continue
		}
		result = append(result, e)
		seen[k] = true
	}

	for _, e := range entries {
		if !seen[e.Key] {
			result = append(result, e)
		}
	}

	return result, nil
}

// ReorderFile reads a .env file, reorders its keys according to order, and
// writes the result back to the same path.
func ReorderFile(path string, order []string, opts ...ReorderOption) error {
	entries, err := ParseFile(path)
	if err != nil {
		return err
	}

	reordered, err := Reorder(entries, order, opts...)
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := WriteFile(tmp, reordered); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

package envfile

import (
	"fmt"
	"os"
)

// CloneOptions controls the behaviour of Clone.
type CloneOptions struct {
	// Overwrite allows the destination file to be replaced if it already exists.
	Overwrite bool
	// Keys restricts cloning to only the specified keys. When empty all keys are copied.
	Keys []string
	// StripValues replaces every value with an empty string in the destination.
	StripValues bool
}

// Clone copies entries from src into a new env map, optionally filtering and
// stripping values. It does not touch the filesystem.
func Clone(entries []Entry, opts CloneOptions) []Entry {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if len(keySet) > 0 {
			if _, ok := keySet[e.Key]; !ok {
				continue
			}
		}
		cloned := Entry{Key: e.Key, Value: e.Value, Comment: e.Comment}
		if opts.StripValues {
			cloned.Value = ""
		}
		out = append(out, cloned)
	}
	return out
}

// CloneFile reads srcPath, applies Clone with opts, and writes the result to
// dstPath. It returns an error if dstPath already exists and Overwrite is false.
func CloneFile(srcPath, dstPath string, opts CloneOptions) error {
	if !opts.Overwrite {
		if _, err := os.Stat(dstPath); err == nil {
			return fmt.Errorf("clone: destination %q already exists (use overwrite to replace)", dstPath)
		}
	}

	entries, err := ParseFile(srcPath)
	if err != nil {
		return fmt.Errorf("clone: read source: %w", err)
	}

	cloned := Clone(entries, opts)

	if err := WriteFile(dstPath, cloned); err != nil {
		return fmt.Errorf("clone: write destination: %w", err)
	}
	return nil
}

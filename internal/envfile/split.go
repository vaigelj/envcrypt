package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SplitOption configures a Split operation.
type SplitOption func(*splitConfig)

type splitConfig struct {
	prefixSep string
	overwrite bool
}

// WithSplitSeparator sets the separator used to detect prefix groups.
func WithSplitSeparator(sep string) SplitOption {
	return func(c *splitConfig) { c.prefixSep = sep }
}

// WithSplitOverwrite allows overwriting existing output files.
func WithSplitOverwrite() SplitOption {
	return func(c *splitConfig) { c.overwrite = true }
}

// Split partitions entries into groups by key prefix and returns a map of
// prefix -> entries. Keys retain their original names.
func Split(entries []Entry, opts ...SplitOption) map[string][]Entry {
	cfg := &splitConfig{prefixSep: "_"}
	for _, o := range opts {
		o(cfg)
	}

	groups := make(map[string][]Entry)
	for _, e := range entries {
		parts := strings.SplitN(e.Key, cfg.prefixSep, 2)
		prefix := parts[0]
		if len(parts) == 1 {
			prefix = "_default"
		}
		groups[prefix] = append(groups[prefix], e)
	}
	return groups
}

// SplitFile reads src, splits by prefix, and writes one file per group into
// outDir named <prefix>.env. Returns the list of files written.
func SplitFile(src, outDir string, opts ...SplitOption) ([]string, error) {
	entries, err := ParseFile(src)
	if err != nil {
		return nil, fmt.Errorf("split: parse %s: %w", src, err)
	}

	cfg := &splitConfig{prefixSep: "_"}
	for _, o := range opts {
		o(cfg)
	}

	groups := Split(entries, opts...)

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("split: mkdir %s: %w", outDir, err)
	}

	var written []string
	for prefix, grp := range groups {
		dest := filepath.Join(outDir, prefix+".env")
		if !cfg.overwrite {
			if _, err := os.Stat(dest); err == nil {
				return written, fmt.Errorf("split: %s already exists (use --overwrite)", dest)
			}
		}
		if err := WriteFile(dest, grp); err != nil {
			return written, fmt.Errorf("split: write %s: %w", dest, err)
		}
		written = append(written, dest)
	}
	return written, nil
}

package envfile

import "fmt"

// PromoteOption configures a Promote operation.
type PromoteOption func(*promoteConfig)

type promoteConfig struct {
	overwrite bool
	exclude   map[string]bool
}

// WithPromoteOverwrite allows existing keys in the target to be overwritten.
func WithPromoteOverwrite() PromoteOption {
	return func(c *promoteConfig) { c.overwrite = true }
}

// WithPromoteExclude skips the given keys during promotion.
func WithPromoteExclude(keys ...string) PromoteOption {
	return func(c *promoteConfig) {
		for _, k := range keys {
			c.exclude[k] = true
		}
	}
}

// PromoteResult holds a summary of what was promoted.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	Conflict []string
}

// Promote copies entries from src into dst according to the given options.
// Keys present in dst are treated as conflicts unless WithPromoteOverwrite is set.
func Promote(src, dst []Entry, opts ...PromoteOption) ([]Entry, PromoteResult, error) {
	cfg := &promoteConfig{exclude: make(map[string]bool)}
	for _, o := range opts {
		o(cfg)
	}

	dstMap := make(map[string]int, len(dst))
	for i, e := range dst {
		dstMap[e.Key] = i
	}

	out := make([]Entry, len(dst))
	copy(out, dst)

	var result PromoteResult
	for _, e := range src {
		if cfg.exclude[e.Key] {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}
		if idx, exists := dstMap[e.Key]; exists {
			if !cfg.overwrite {
				result.Conflict = append(result.Conflict, e.Key)
				continue
			}
			out[idx] = e
			result.Promoted = append(result.Promoted, e.Key)
		} else {
			out = append(out, e)
			dstMap[e.Key] = len(out) - 1
			result.Promoted = append(result.Promoted, e.Key)
		}
	}
	return out, result, nil
}

// PromoteFile reads src and dst files, promotes entries, and writes the result
// back to the dst path.
func PromoteFile(srcPath, dstPath string, opts ...PromoteOption) (PromoteResult, error) {
	src, err := ParseFile(srcPath)
	if err != nil {
		return PromoteResult{}, fmt.Errorf("promote: read src: %w", err)
	}
	dst, err := ParseFile(dstPath)
	if err != nil {
		return PromoteResult{}, fmt.Errorf("promote: read dst: %w", err)
	}
	out, result, err := Promote(src, dst, opts...)
	if err != nil {
		return result, err
	}
	if err := WriteFile(dstPath, out); err != nil {
		return result, fmt.Errorf("promote: write dst: %w", err)
	}
	return result, nil
}

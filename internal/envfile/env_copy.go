package envfile

import "fmt"

// CopyOptions controls how CopyEnv behaves.
type CopyOptions struct {
	// Overwrite existing keys in dst when true.
	Overwrite bool
	// Keys to exclude from the copy.
	Exclude []string
}

// CopyEnv copies key/value pairs from src into dst according to opts.
// It returns the number of keys copied.
func CopyEnv(dst, src map[string]string, opts CopyOptions) int {
	excluded := make(map[string]bool, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excluded[k] = true
	}

	copied := 0
	for k, v := range src {
		if excluded[k] {
			continue
		}
		if _, exists := dst[k]; exists && !opts.Overwrite {
			continue
		}
		dst[k] = v
		copied++
	}
	return copied
}

// CopyFile reads src path, copies entries into dst path and writes the result.
func CopyFile(dstPath, srcPath string, opts CopyOptions) (int, error) {
	src, err := ParseFile(srcPath)
	if err != nil {
		return 0, fmt.Errorf("reading src: %w", err)
	}

	dst, err := ParseFile(dstPath)
	if err != nil {
		// dst may not exist yet — start empty
		dst = map[string]string{}
	}

	n := CopyEnv(dst, src, opts)
	if err := WriteFile(dstPath, dst); err != nil {
		return 0, fmt.Errorf("writing dst: %w", err)
	}
	return n, nil
}

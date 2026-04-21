package envfile

import "fmt"

// Chain represents an ordered list of .env file paths whose entries are
// merged left-to-right, with later files taking precedence over earlier ones.
type Chain struct {
	Files []string
}

// NewChain creates a Chain from an ordered list of file paths.
func NewChain(files ...string) *Chain {
	return &Chain{Files: files}
}

// Resolve loads and merges all files in the chain.
// Entries from later files override entries from earlier files.
func (c *Chain) Resolve() ([]Entry, error) {
	if len(c.Files) == 0 {
		return nil, nil
	}

	base, err := ParseFile(c.Files[0])
	if err != nil {
		return nil, fmt.Errorf("chain: loading %q: %w", c.Files[0], err)
	}

	for _, path := range c.Files[1:] {
		override, err := ParseFile(path)
		if err != nil {
			return nil, fmt.Errorf("chain: loading %q: %w", path, err)
		}
		base = Merge(base, override, true)
	}

	return base, nil
}

// ResolveMap is a convenience wrapper that returns the merged entries as a
// key→value map.
func (c *Chain) ResolveMap() (map[string]string, error) {
	entries, err := c.Resolve()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(entries))
	for _, e := range entries {
		out[e.Key] = e.Value
	}
	return out, nil
}

// Sources returns the list of file paths in the chain.
func (c *Chain) Sources() []string {
	result := make([]string, len(c.Files))
	copy(result, c.Files)
	return result
}

package envfile

// DedupeOption controls how duplicate keys are resolved.
type DedupeOption int

const (
	// DedupeKeepFirst retains the first occurrence of a duplicate key.
	DedupeKeepFirst DedupeOption = iota
	// DedupeKeepLast retains the last occurrence of a duplicate key.
	DedupeKeepLast
)

// DedupeResult holds the outcome of a deduplication pass.
type DedupeResult struct {
	Entries    []Entry
	Duplicates []string // keys that had duplicates removed
}

// Dedupe removes duplicate keys from a slice of entries according to opt.
// Order of surviving entries is preserved (relative to the chosen occurrence).
func Dedupe(entries []Entry, opt DedupeOption) DedupeResult {
	seen := make(map[string]int) // key -> index in out
	out := make([]Entry, 0, len(entries))
	dupeSet := make(map[string]bool)

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			dupeSet[e.Key] = true
			if opt == DedupeKeepLast {
				out[idx] = e
			}
			// DedupeKeepFirst: do nothing, keep original
		} else {
			seen[e.Key] = len(out)
			out = append(out, e)
		}
	}

	dupes := make([]string, 0, len(dupeSet))
	for k := range dupeSet {
		dupes = append(dupes, k)
	}

	return DedupeResult{Entries: out, Duplicates: dupes}
}

// DedupeFile reads the file at path, deduplicates its entries using opt,
// and writes the result back to the same file.
func DedupeFile(path string, opt DedupeOption) (DedupeResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return DedupeResult{}, err
	}
	res := Dedupe(entries, opt)
	if err := WriteFile(path, res.Entries); err != nil {
		return DedupeResult{}, err
	}
	return res, nil
}

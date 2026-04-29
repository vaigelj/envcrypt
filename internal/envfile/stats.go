package envfile

import "sort"

// Stats holds aggregate statistics about a set of env entries.
type Stats struct {
	Total        int
	WithValues   int
	Empty        int
	Comments     int
	Unique       int
	Duplicates   int
	LongestKey   string
	LongestValue string
	TopKeys      []string // top 5 keys by value length
}

// GatherStats computes statistics over the provided entries.
func GatherStats(entries []Entry) Stats {
	seen := make(map[string]int)
	var s Stats

	for _, e := range entries {
		if e.Key == "" {
			s.Comments++
			continue
		}
		s.Total++
		seen[e.Key]++
		if e.Value == "" {
			s.Empty++
		} else {
			s.WithValues++
		}
		if len(e.Key) > len(s.LongestKey) {
			s.LongestKey = e.Key
		}
		if len(e.Value) > len(s.LongestValue) {
			s.LongestValue = e.Value
		}
	}

	for k, count := range seen {
		if count == 1 {
			s.Unique++
		} else {
			s.Duplicates += count - 1
			_ = k
		}
	}

	// Build top-5 keys by value length.
	type kv struct {
		key string
		len int
	}
	var pairs []kv
	// Deduplicate for top-keys: use last seen value length per key.
	keyLen := make(map[string]int)
	for _, e := range entries {
		if e.Key != "" {
			keyLen[e.Key] = len(e.Value)
		}
	}
	for k, l := range keyLen {
		pairs = append(pairs, kv{k, l})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].len != pairs[j].len {
			return pairs[i].len > pairs[j].len
		}
		return pairs[i].key < pairs[j].key
	})
	for i := 0; i < len(pairs) && i < 5; i++ {
		s.TopKeys = append(s.TopKeys, pairs[i].key)
	}

	return s
}

// GatherStatsFile reads a file and returns its stats.
func GatherStatsFile(path string) (Stats, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return Stats{}, err
	}
	return GatherStats(entries), nil
}

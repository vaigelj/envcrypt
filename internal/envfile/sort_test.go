package envfile

import (
	"testing"
)

func makeEntries(keys ...string) []Entry {
	out := make([]Entry, len(keys))
	for i, k := range keys {
		out[i] = Entry{Key: k, Value: "v"}
	}
	return out
}

func keys(entries []Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Key
	}
	return out
}

func TestSortAlpha(t *testing.T) {
	entries := makeEntries("ZEBRA", "APPLE", "MANGO")
	sorted := Sort(entries, SortOptions{Order: SortAlpha})
	got := keys(sorted)
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("pos %d: got %s want %s", i, got[i], k)
		}
	}
}

func TestSortAlphaDesc(t *testing.T) {
	entries := makeEntries("APPLE", "ZEBRA", "MANGO")
	sorted := Sort(entries, SortOptions{Order: SortAlphaDesc})
	got := keys(sorted)
	want := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("pos %d: got %s want %s", i, got[i], k)
		}
	}
}

func TestSortByLength(t *testing.T) {
	entries := makeEntries("LONGKEY", "A", "MED")
	sorted := Sort(entries, SortOptions{Order: SortByLength})
	got := keys(sorted)
	want := []string{"A", "MED", "LONGKEY"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("pos %d: got %s want %s", i, got[i], k)
		}
	}
}

func TestSortWithGroups(t *testing.T) {
	entries := makeEntries("ZEBRA", "DB_HOST", "APPLE", "DB_PORT")
	sorted := Sort(entries, SortOptions{Order: SortAlpha, Groups: []string{"DB_"}})
	got := keys(sorted)
	// DB_ entries should come first, then rest alphabetically
	if got[0] != "DB_HOST" || got[1] != "DB_PORT" {
		t.Errorf("expected DB_ entries first, got %v", got)
	}
	if got[2] != "APPLE" || got[3] != "ZEBRA" {
		t.Errorf("expected remaining alpha order, got %v", got[2:])
	}
}

func TestSortPreservesOriginal(t *testing.T) {
	entries := makeEntries("B", "A")
	_ = Sort(entries, SortOptions{Order: SortAlpha})
	if entries[0].Key != "B" {
		t.Error("Sort should not modify original slice")
	}
}

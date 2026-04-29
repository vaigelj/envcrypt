package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// SummaryReport holds statistics about a parsed env file.
type SummaryReport struct {
	TotalKeys    int
	EmptyValues  int
	CommentLines int
	UniqueKeys   int
	DuplicateKeys []string
	SensitiveKeys []string
	LongestKey   string
	LongestValue string
}

// String returns a human-readable summary report.
func (r SummaryReport) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Total keys     : %d\n", r.TotalKeys)
	fmt.Fprintf(&sb, "Unique keys    : %d\n", r.UniqueKeys)
	fmt.Fprintf(&sb, "Empty values   : %d\n", r.EmptyValues)
	fmt.Fprintf(&sb, "Comment lines  : %d\n", r.CommentLines)
	if len(r.DuplicateKeys) > 0 {
		fmt.Fprintf(&sb, "Duplicate keys : %s\n", strings.Join(r.DuplicateKeys, ", "))
	}
	if len(r.SensitiveKeys) > 0 {
		fmt.Fprintf(&sb, "Sensitive keys : %s\n", strings.Join(r.SensitiveKeys, ", "))
	}
	if r.LongestKey != "" {
		fmt.Fprintf(&sb, "Longest key    : %s\n", r.LongestKey)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Summarize computes statistics for a slice of Entry values.
func Summarize(entries []Entry) SummaryReport {
	var report SummaryReport
	seen := make(map[string]int)

	for _, e := range entries {
		if e.Comment {
			report.CommentLines++
			continue
		}
		report.TotalKeys++
		seen[e.Key]++
		if strings.TrimSpace(e.Value) == "" {
			report.EmptyValues++
		}
		if IsSensitive(e.Key) {
			report.SensitiveKeys = append(report.SensitiveKeys, e.Key)
		}
		if len(e.Key) > len(report.LongestKey) {
			report.LongestKey = e.Key
		}
		if len(e.Value) > len(report.LongestValue) {
			report.LongestValue = e.Key // store key name of longest value
		}
	}

	dupes := []string{}
	for k, count := range seen {
		if count > 1 {
			dupes = append(dupes, k)
		}
	}
	sort.Strings(dupes)
	report.UniqueKeys = len(seen)
	report.DuplicateKeys = dupes

	return report
}

// SummarizeFile parses the file at path and returns a SummaryReport.
func SummarizeFile(path string) (SummaryReport, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return SummaryReport{}, fmt.Errorf("summarize: %w", err)
	}
	return Summarize(entries), nil
}

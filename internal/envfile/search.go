package envfile

import (
	"regexp"
	"strings"
)

// SearchResult holds a matched key-value pair and the file it came from.
type SearchResult struct {
	Key   string
	Value string
	File  string
}

// SearchOptions controls how Search behaves.
type SearchOptions struct {
	// KeyPattern filters by key using a substring or regex if UseRegex is true.
	KeyPattern string
	// ValuePattern filters by value using a substring or regex if UseRegex is true.
	ValuePattern string
	UseRegex     bool
	CaseSensitive bool
}

// Search scans env entries and returns those matching the given options.
func Search(entries []Entry, file string, opts SearchOptions) ([]SearchResult, error) {
	keyRe, err := compilePattern(opts.KeyPattern, opts.UseRegex, opts.CaseSensitive)
	if err != nil {
		return nil, err
	}
	valRe, err := compilePattern(opts.ValuePattern, opts.UseRegex, opts.CaseSensitive)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, e := range entries {
		if !matchPattern(keyRe, opts.KeyPattern, e.Key, opts.CaseSensitive) {
			continue
		}
		if !matchPattern(valRe, opts.ValuePattern, e.Value, opts.CaseSensitive) {
			continue
		}
		results = append(results, SearchResult{Key: e.Key, Value: e.Value, File: file})
	}
	return results, nil
}

// SearchFile parses the given file and runs Search on it.
func SearchFile(path string, opts SearchOptions) ([]SearchResult, error) {
	entries, err := ParseFile(path)
	if err != nil {
		return nil, err
	}
	return Search(entries, path, opts)
}

func compilePattern(pattern string, useRegex, caseSensitive bool) (*regexp.Regexp, error) {
	if pattern == "" || !useRegex {
		return nil, nil
	}
	p := pattern
	if !caseSensitive {
		p = "(?i)" + p
	}
	return regexp.Compile(p)
}

func matchPattern(re *regexp.Regexp, pattern, value string, caseSensitive bool) bool {
	if pattern == "" {
		return true
	}
	if re != nil {
		return re.MatchString(value)
	}
	if !caseSensitive {
		return strings.Contains(strings.ToLower(value), strings.ToLower(pattern))
	}
	return strings.Contains(value, pattern)
}

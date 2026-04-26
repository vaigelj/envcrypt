// Package envfile provides utilities for parsing, writing, and manipulating .env files.
package envfile

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ComputeOp represents a supported arithmetic or string operation.
type ComputeOp string

const (
	ComputeAdd    ComputeOp = "add"
	ComputeSub    ComputeOp = "sub"
	ComputeMul    ComputeOp = "mul"
	ComputeDiv    ComputeOp = "div"
	ComputeConcat ComputeOp = "concat"
	ComputeUpper  ComputeOp = "upper"
	ComputeLower  ComputeOp = "lower"
)

// ComputeRule describes a single derived-value rule: read from Source keys,
// apply Op with optional Args, and write the result to Dest.
type ComputeRule struct {
	Dest   string    // key to write the result into
	Op     ComputeOp // operation to perform
	Source []string  // source keys (resolved from entries)
	Args   []string  // extra literal arguments (e.g. separator for concat)
}

// Compute applies a set of ComputeRules to entries, adding or updating the
// destination key with the computed value. Source keys are resolved from the
// current entry map; missing keys fall back to the process environment.
// Returns the modified slice and any computation error.
func Compute(entries []Entry, rules []ComputeRule) ([]Entry, error) {
	// Build a lookup map for fast access.
	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	resolve := func(key string) (string, bool) {
		if i, ok := index[key]; ok {
			return entries[i].Value, true
		}
		if v, ok := os.LookupEnv(key); ok {
			return v, true
		}
		return "", false
	}

	for _, rule := range rules {
		result, err := applyComputeOp(rule, resolve)
		if err != nil {
			return entries, fmt.Errorf("compute rule %q: %w", rule.Dest, err)
		}

		if i, ok := index[rule.Dest]; ok {
			entries[i].Value = result
		} else {
			entries = append(entries, Entry{Key: rule.Dest, Value: result})
			index[rule.Dest] = len(entries) - 1
		}
	}
	return entries, nil
}

// ComputeFile reads entries from src, applies rules, and writes the result to
// dst. If dst is empty the source file is overwritten.
func ComputeFile(src, dst string, rules []ComputeRule) error {
	entries, err := ParseFile(src)
	if err != nil {
		return err
	}
	entries, err = Compute(entries, rules)
	if err != nil {
		return err
	}
	if dst == "" {
		dst = src
	}
	return WriteFile(dst, entries)
}

// applyComputeOp executes a single rule and returns the string result.
func applyComputeOp(rule ComputeRule, resolve func(string) (string, bool)) (string, error) {
	switch rule.Op {
	case ComputeUpper:
		if len(rule.Source) != 1 {
			return "", fmt.Errorf("upper requires exactly 1 source key")
		}
		v, _ := resolve(rule.Source[0])
		return strings.ToUpper(v), nil

	case ComputeLower:
		if len(rule.Source) != 1 {
			return "", fmt.Errorf("lower requires exactly 1 source key")
		}
		v, _ := resolve(rule.Source[0])
		return strings.ToLower(v), nil

	case ComputeConcat:
		sep := ""
		if len(rule.Args) > 0 {
			sep = rule.Args[0]
		}
		parts := make([]string, 0, len(rule.Source))
		for _, k := range rule.Source {
			v, _ := resolve(k)
			parts = append(parts, v)
		}
		return strings.Join(parts, sep), nil

	case ComputeAdd, ComputeSub, ComputeMul, ComputeDiv:
		if len(rule.Source) < 2 {
			return "", fmt.Errorf("%s requires at least 2 source keys", rule.Op)
		}
		first, ok := resolve(rule.Source[0])
		if !ok {
			return "", fmt.Errorf("source key %q not found", rule.Source[0])
		}
		acc, err := strconv.ParseFloat(first, 64)
		if err != nil {
			return "", fmt.Errorf("source key %q is not numeric: %w", rule.Source[0], err)
		}
		for _, k := range rule.Source[1:] {
			raw, ok := resolve(k)
			if !ok {
				return "", fmt.Errorf("source key %q not found", k)
			}
			v, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return "", fmt.Errorf("source key %q is not numeric: %w", k, err)
			}
			switch rule.Op {
			case ComputeAdd:
				acc += v
			case ComputeSub:
				acc -= v
			case ComputeMul:
				acc *= v
			case ComputeDiv:
				if v == 0 {
					return "", fmt.Errorf("division by zero (key %q)", k)
				}
				acc /= v
			}
		}
		// Format as integer when result is whole, otherwise use float.
		if acc == float64(int64(acc)) {
			return strconv.FormatInt(int64(acc), 10), nil
		}
		return strconv.FormatFloat(acc, 'f', -1, 64), nil

	default:
		return "", fmt.Errorf("unknown operation %q", rule.Op)
	}
}

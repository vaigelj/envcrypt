package envfile

import (
	"fmt"
	"strings"
)

// MaskMode controls how values are masked in output.
type MaskMode int

const (
	// MaskFull replaces the entire value with asterisks.
	MaskFull MaskMode = iota
	// MaskPartial reveals the first and last characters.
	MaskPartial
	// MaskLength replaces the value with a fixed-length mask.
	MaskLength
)

// MaskOptions configures masking behaviour.
type MaskOptions struct {
	Mode       MaskMode
	MaskChar   rune
	FixedLen   int  // used by MaskLength
	RevealLen  int  // used by MaskPartial (chars revealed on each side)
	Keys       []string // if non-empty, only mask these keys; otherwise mask all sensitive keys
}

var defaultMaskOptions = MaskOptions{
	Mode:      MaskFull,
	MaskChar:  '*',
	FixedLen:  8,
	RevealLen: 2,
}

// MaskValue applies the masking strategy to a single value string.
func MaskValue(value string, opts MaskOptions) string {
	if value == "" {
		return ""
	}
	ch := string(opts.MaskChar)
	switch opts.Mode {
	case MaskPartial:
		n := opts.RevealLen
		if n <= 0 {
			n = 2
		}
		runes := []rune(value)
		if len(runes) <= n*2 {
			return strings.Repeat(ch, len(runes))
		}
		return string(runes[:n]) + strings.Repeat(ch, len(runes)-n*2) + string(runes[len(runes)-n:])
	case MaskLength:
		l := opts.FixedLen
		if l <= 0 {
			l = 8
		}
		return strings.Repeat(ch, l)
	default: // MaskFull
		return strings.Repeat(ch, len([]rune(value)))
	}
}

// MaskEntries returns a new slice of Entry with sensitive (or specified) values masked.
func MaskEntries(entries []Entry, opts MaskOptions) []Entry {
	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = true
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		copy := e
		shouldMask := len(keySet) > 0 && keySet[strings.ToUpper(e.Key)] ||
			len(keySet) == 0 && IsSensitive(e.Key)
		if shouldMask {
			copy.Value = MaskValue(e.Value, opts)
		}
		result[i] = copy
	}
	return result
}

// MaskSummary returns a human-readable summary of which keys were masked.
func MaskSummary(original, masked []Entry) string {
	var sb strings.Builder
	for i, o := range original {
		if i < len(masked) && masked[i].Value != o.Value {
			fmt.Fprintf(&sb, "  %s: masked (%d chars)\n", o.Key, len([]rune(o.Value)))
		}
	}
	return sb.String()
}

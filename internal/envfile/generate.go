package envfile

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	charsetAlpha   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetNumeric = "0123456789"
	charsetSymbol  = "!@#$%^&*()-_=+[]{}"
	charsetAll     = charsetAlpha + charsetNumeric + charsetSymbol
)

// GenerateOptions controls value generation.
type GenerateOptions struct {
	Length  int
	NoSymbols bool
	Numeric   bool
}

// GenerateValue returns a cryptographically random string.
func GenerateValue(opts GenerateOptions) (string, error) {
	if opts.Length <= 0 {
		opts.Length = 32
	}
	charset := charsetAll
	if opts.Numeric {
		charset = charsetNumeric
	} else if opts.NoSymbols {
		charset = charsetAlpha + charsetNumeric
	}
	var sb strings.Builder
	for i := 0; i < opts.Length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("generate: %w", err)
		}
		sb.WriteByte(charset[idx.Int64()])
	}
	return sb.String(), nil
}

// GenerateForKeys returns a map of key -> generated value for the given keys.
func GenerateForKeys(keys []string, opts GenerateOptions) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		v, err := GenerateValue(opts)
		if err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}

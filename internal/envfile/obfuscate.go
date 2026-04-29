package envfile

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

// ObfuscateOption configures obfuscation behaviour.
type ObfuscateOption func(*obfuscateConfig)

type obfuscateConfig struct {
	keys    map[string]bool
	prefix  string
	padding int
}

// WithObfuscateKeys limits obfuscation to the given keys.
func WithObfuscateKeys(keys []string) ObfuscateOption {
	return func(c *obfuscateConfig) {
		for _, k := range keys {
			c.keys[k] = true
		}
	}
}

// WithObfuscatePrefix sets a prefix that marks obfuscated values.
func WithObfuscatePrefix(prefix string) ObfuscateOption {
	return func(c *obfuscateConfig) {
		c.prefix = prefix
	}
}

// WithObfuscatePadding adds random padding bytes to each obfuscated value.
func WithObfuscatePadding(n int) ObfuscateOption {
	return func(c *obfuscateConfig) {
		c.padding = n
	}
}

const defaultObfuscatePrefix = "obf:"

// Obfuscate base64-encodes entry values (optionally with random padding)
// so they are not immediately human-readable. It is not encryption;
// use EncryptFields for cryptographic protection.
func Obfuscate(entries []Entry, opts ...ObfuscateOption) ([]Entry, error) {
	cfg := &obfuscateConfig{
		keys:   make(map[string]bool),
		prefix: defaultObfuscatePrefix,
	}
	for _, o := range opts {
		o(cfg)
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Comment || e.Blank {
			out[i] = e
			continue
		}
		if len(cfg.keys) > 0 && !cfg.keys[e.Key] {
			out[i] = e
			continue
		}
		if strings.HasPrefix(e.Value, cfg.prefix) {
			out[i] = e
			continue
		}
		raw := []byte(e.Value)
		if cfg.padding > 0 {
			pad := make([]byte, cfg.padding)
			if _, err := rand.Read(pad); err != nil {
				return nil, fmt.Errorf("obfuscate: generate padding: %w", err)
			}
			raw = append(pad, raw...)
		}
		encoded := cfg.prefix + base64.StdEncoding.EncodeToString(raw)
		out[i] = Entry{Key: e.Key, Value: encoded, RawLine: e.RawLine}
	}
	return out, nil
}

// Deobfuscate reverses Obfuscate, stripping the prefix and decoding base64.
func Deobfuscate(entries []Entry, opts ...ObfuscateOption) ([]Entry, error) {
	cfg := &obfuscateConfig{
		keys:   make(map[string]bool),
		prefix: defaultObfuscatePrefix,
		padding: 0,
	}
	for _, o := range opts {
		o(cfg)
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Comment || e.Blank {
			out[i] = e
			continue
		}
		if len(cfg.keys) > 0 && !cfg.keys[e.Key] {
			out[i] = e
			continue
		}
		if !strings.HasPrefix(e.Value, cfg.prefix) {
			out[i] = e
			continue
		}
		b64 := strings.TrimPrefix(e.Value, cfg.prefix)
		raw, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, fmt.Errorf("deobfuscate key %q: %w", e.Key, err)
		}
		if cfg.padding > 0 && len(raw) >= cfg.padding {
			raw = raw[cfg.padding:]
		}
		out[i] = Entry{Key: e.Key, Value: string(raw), RawLine: e.RawLine}
	}
	return out, nil
}

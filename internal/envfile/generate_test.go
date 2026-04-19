package envfile

import (
	"strings"
	"testing"
)

func TestGenerateValueDefaultLength(t *testing.T) {
	v, err := GenerateValue(GenerateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 32 {
		t.Fatalf("expected length 32, got %d", len(v))
	}
}

func TestGenerateValueCustomLength(t *testing.T) {
	v, err := GenerateValue(GenerateOptions{Length: 16})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 16 {
		t.Fatalf("expected length 16, got %d", len(v))
	}
}

func TestGenerateValueNumeric(t *testing.T) {
	v, err := GenerateValue(GenerateOptions{Length: 20, Numeric: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range v {
		if !strings.ContainsRune(charsetNumeric, c) {
			t.Fatalf("non-numeric character %q in output", c)
		}
	}
}

func TestGenerateValueNoSymbols(t *testing.T) {
	v, err := GenerateValue(GenerateOptions{Length: 40, NoSymbols: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range v {
		if strings.ContainsRune(charsetSymbol, c) {
			t.Fatalf("symbol character %q found in no-symbols output", c)
		}
	}
}

func TestGenerateForKeys(t *testing.T) {
	keys := []string{"SECRET_A", "SECRET_B", "TOKEN"}
	out, err := GenerateForKeys(keys, GenerateOptions{Length: 24})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != len(keys) {
		t.Fatalf("expected %d entries, got %d", len(keys), len(out))
	}
	for _, k := range keys {
		if len(out[k]) != 24 {
			t.Fatalf("key %s: expected length 24, got %d", k, len(out[k]))
		}
	}
	// values should be unique
	if out["SECRET_A"] == out["SECRET_B"] {
		t.Fatal("expected unique values for different keys")
	}
}

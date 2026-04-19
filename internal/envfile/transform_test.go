package envfile

import (
	"testing"
)

func TestTransformAll(t *testing.T) {
	pairs := []Entry{{Key: "FOO", Value: "hello"}, {Key: "BAR", Value: "world"}}
	out := Transform(pairs, UppercaseValues(), TransformOptions{})
	if out[0].Value != "HELLO" || out[1].Value != "WORLD" {
		t.Fatalf("unexpected values: %v", out)
	}
}

func TestTransformSpecificKeys(t *testing.T) {
	pairs := []Entry{{Key: "FOO", Value: "hello"}, {Key: "BAR", Value: "world"}}
	out := Transform(pairs, UppercaseValues(), TransformOptions{Keys: []string{"FOO"}})
	if out[0].Value != "HELLO" {
		t.Fatalf("expected FOO uppercased, got %s", out[0].Value)
	}
	if out[1].Value != "world" {
		t.Fatalf("expected BAR unchanged, got %s", out[1].Value)
	}
}

func TestTransformExclude(t *testing.T) {
	pairs := []Entry{{Key: "FOO", Value: "hello"}, {Key: "BAR", Value: "world"}}
	out := Transform(pairs, UppercaseValues(), TransformOptions{Exclude: []string{"BAR"}})
	if out[0].Value != "HELLO" {
		t.Fatalf("expected FOO uppercased, got %s", out[0].Value)
	}
	if out[1].Value != "world" {
		t.Fatalf("expected BAR excluded, got %s", out[1].Value)
	}
}

func TestTrimValues(t *testing.T) {
	pairs := []Entry{{Key: "A", Value: "  spaced  "}}
	out := Transform(pairs, TrimValues(), TransformOptions{})
	if out[0].Value != "spaced" {
		t.Fatalf("expected trimmed value, got %q", out[0].Value)
	}
}

func TestPrefixValues(t *testing.T) {
	pairs := []Entry{{Key: "URL", Value: "example.com"}}
	out := Transform(pairs, PrefixValues("https://"), TransformOptions{})
	if out[0].Value != "https://example.com" {
		t.Fatalf("unexpected value: %s", out[0].Value)
	}
}

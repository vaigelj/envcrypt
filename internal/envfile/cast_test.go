package envfile

import (
	"testing"
)

func makeCastEntries() []Entry {
	return []Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "RATIO", Value: "3.14"},
		{Key: "ENABLED", Value: "yes"},
		{Key: "LABEL", Value: "  hello  "},
	}
}

func TestCastToInt(t *testing.T) {
	entries, err := Cast(makeCastEntries(), WithCastKeys(map[string]string{"PORT": "int"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := entries[0].Value; got != "8080" {
		t.Errorf("PORT: got %q, want %q", got, "8080")
	}
}

func TestCastFloatToInt(t *testing.T) {
	entries, err := Cast(makeCastEntries(), WithCastKeys(map[string]string{"RATIO": "int"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := entries[1].Value; got != "3" {
		t.Errorf("RATIO as int: got %q, want %q", got, "3")
	}
}

func TestCastToBool(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"yes", "true"}, {"1", "true"}, {"on", "true"}, {"true", "true"},
		{"no", "false"}, {"0", "false"}, {"off", "false"}, {"false", "false"},
	}
	for _, tc := range cases {
		entries := []Entry{{Key: "FLAG", Value: tc.input}}
		out, err := Cast(entries, WithCastKeys(map[string]string{"FLAG": "bool"}))
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if out[0].Value != tc.want {
			t.Errorf("input %q: got %q, want %q", tc.input, out[0].Value, tc.want)
		}
	}
}

func TestCastInvalidBoolLenient(t *testing.T) {
	entries := []Entry{{Key: "FLAG", Value: "maybe"}}
	out, err := Cast(entries, WithCastKeys(map[string]string{"FLAG": "bool"}))
	if err != nil {
		t.Fatalf("unexpected error in lenient mode: %v", err)
	}
	// value unchanged
	if out[0].Value != "maybe" {
		t.Errorf("expected value unchanged, got %q", out[0].Value)
	}
}

func TestCastInvalidBoolStrict(t *testing.T) {
	entries := []Entry{{Key: "FLAG", Value: "maybe"}}
	_, err := Cast(entries,
		WithCastKeys(map[string]string{"FLAG": "bool"}),
		WithCastStrict(),
	)
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestCastUnknownKeyUnchanged(t *testing.T) {
	entries := makeCastEntries()
	out, err := Cast(entries, WithCastKeys(map[string]string{"NONEXISTENT": "int"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range out {
		if e.Value != entries[i].Value {
			t.Errorf("key %q: value changed unexpectedly", e.Key)
		}
	}
}

func TestCastUnknownType(t *testing.T) {
	entries := []Entry{{Key: "X", Value: "foo"}}
	_, err := Cast(entries,
		WithCastKeys(map[string]string{"X": "uuid"}),
		WithCastStrict(),
	)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

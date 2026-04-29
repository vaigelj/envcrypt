package envfile

import (
	"strings"
	"testing"
)

func makeObfuscateEntries() []Entry {
	return []Entry{
		{Key: "DB_PASS", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "APP_NAME", Value: "myapp"},
		{Blank: true},
		{Comment: true, RawLine: "# comment"},
	}
}

func TestObfuscateAll(t *testing.T) {
	entries := makeObfuscateEntries()
	out, err := Obfuscate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Blank || e.Comment {
			continue
		}
		if !strings.HasPrefix(e.Value, defaultObfuscatePrefix) {
			t.Errorf("key %q value %q missing prefix", e.Key, e.Value)
		}
	}
}

func TestObfuscateSubset(t *testing.T) {
	entries := makeObfuscateEntries()
	out, err := Obfuscate(entries, WithObfuscateKeys([]string{"DB_PASS"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Blank || e.Comment {
			continue
		}
		if e.Key == "DB_PASS" {
			if !strings.HasPrefix(e.Value, defaultObfuscatePrefix) {
				t.Errorf("DB_PASS should be obfuscated")
			}
		} else {
			if strings.HasPrefix(e.Value, defaultObfuscatePrefix) {
				t.Errorf("key %q should not be obfuscated", e.Key)
			}
		}
	}
}

func TestDeobfuscateRoundtrip(t *testing.T) {
	entries := makeObfuscateEntries()
	obf, err := Obfuscate(entries)
	if err != nil {
		t.Fatalf("obfuscate: %v", err)
	}
	got, err := Deobfuscate(obf)
	if err != nil {
		t.Fatalf("deobfuscate: %v", err)
	}
	for i, orig := range entries {
		if orig.Blank || orig.Comment {
			continue
		}
		if got[i].Value != orig.Value {
			t.Errorf("key %q: got %q, want %q", orig.Key, got[i].Value, orig.Value)
		}
	}
}

func TestObfuscateIdempotent(t *testing.T) {
	entries := makeObfuscateEntries()
	first, _ := Obfuscate(entries)
	second, err := Obfuscate(first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range first {
		if e.Blank || e.Comment {
			continue
		}
		if second[i].Value != e.Value {
			t.Errorf("key %q: double obfuscation changed value", e.Key)
		}
	}
}

func TestObfuscateCustomPrefix(t *testing.T) {
	entries := []Entry{{Key: "SECRET", Value: "hunter2"}}
	out, err := Obfuscate(entries, WithObfuscatePrefix("enc:"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out[0].Value, "enc:") {
		t.Errorf("expected enc: prefix, got %q", out[0].Value)
	}
	back, err := Deobfuscate(out, WithObfuscatePrefix("enc:"))
	if err != nil {
		t.Fatalf("deobfuscate: %v", err)
	}
	if back[0].Value != "hunter2" {
		t.Errorf("got %q, want hunter2", back[0].Value)
	}
}

func TestObfuscateWithPadding(t *testing.T) {
	entries := []Entry{{Key: "TOKEN", Value: "mytoken"}}
	out, err := Obfuscate(entries, WithObfuscatePadding(8))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	back, err := Deobfuscate(out, WithObfuscatePadding(8))
	if err != nil {
		t.Fatalf("deobfuscate: %v", err)
	}
	if back[0].Value != "mytoken" {
		t.Errorf("got %q, want mytoken", back[0].Value)
	}
}

package envfile

import (
	"os"
	"testing"
)

func TestNormalizeUpperKeys(t *testing.T) {
	in := []Entry{{Key: "db_host", Value: "localhost"}, {Key: "api_key", Value: "secret"}}
	out := Normalize(in, WithUpperKeys())
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out[0].Key)
	}
	if out[1].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", out[1].Key)
	}
}

func TestNormalizeTrimValues(t *testing.T) {
	in := []Entry{{Key: "FOO", Value: "  bar  "}, {Key: "BAZ", Value: "\t qux\t"}}
	out := Normalize(in, WithTrimValues())
	if out[0].Value != "bar" {
		t.Errorf("expected 'bar', got %q", out[0].Value)
	}
	if out[1].Value != "qux" {
		t.Errorf("expected 'qux', got %q", out[1].Value)
	}
}

func TestNormalizeRemoveEmpty(t *testing.T) {
	in := []Entry{{Key: "FOO", Value: ""}, {Key: "BAR", Value: "hello"}, {Key: "BAZ", Value: ""}}
	out := Normalize(in, WithTrimValues(), WithRemoveEmpty())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "BAR" {
		t.Errorf("expected BAR, got %s", out[0].Key)
	}
}

func TestNormalizeQuoteValues(t *testing.T) {
	in := []Entry{
		{Key: "GREETING", Value: "hello world"},
		{Key: "TOKEN", Value: "abc123"},
		{Key: "ALREADY", Value: `"already quoted"`},
	}
	out := Normalize(in, WithQuoteValues())
	if out[0].Value != `"hello world"` {
		t.Errorf("expected quoted value, got %q", out[0].Value)
	}
	if out[1].Value != "abc123" {
		t.Errorf("expected unchanged value, got %q", out[1].Value)
	}
	if out[2].Value != `"already quoted"` {
		t.Errorf("expected unchanged already-quoted value, got %q", out[2].Value)
	}
}

func TestNormalizeCombined(t *testing.T) {
	in := []Entry{
		{Key: "db_url", Value: "  postgres://localhost  "},
		{Key: "empty_val", Value: "  "},
	}
	out := Normalize(in, WithUpperKeys(), WithTrimValues(), WithRemoveEmpty())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "DB_URL" || out[0].Value != "postgres://localhost" {
		t.Errorf("unexpected entry: %+v", out[0])
	}
}

func TestNormalizeFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("db_host=  localhost  \napi_key=  secret  \n")
	f.Close()

	if err := NormalizeFile(f.Name(), WithUpperKeys(), WithTrimValues()); err != nil {
		t.Fatalf("NormalizeFile: %v", err)
	}

	entries, err := ParseFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "DB_HOST" || entries[0].Value != "localhost" {
		t.Errorf("unexpected entry[0]: %+v", entries[0])
	}
	if entries[1].Key != "API_KEY" || entries[1].Value != "secret" {
		t.Errorf("unexpected entry[1]: %+v", entries[1])
	}
}

package envfile

import (
	"testing"
)

func TestTokenizeBlankLine(t *testing.T) {
	toks, err := Tokenize([]string{"", "   "})
	if err != nil {
		t.Fatal(err)
	}
	for _, tok := range toks {
		if tok.Type != TokenBlank {
			t.Errorf("expected TokenBlank, got %v", tok.Type)
		}
	}
}

func TestTokenizeComment(t *testing.T) {
	toks, err := Tokenize([]string{"# this is a comment"})
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) != 1 || toks[0].Type != TokenComment {
		t.Fatalf("expected one comment token, got %+v", toks)
	}
	if toks[0].Comment != " this is a comment" {
		t.Errorf("unexpected comment text: %q", toks[0].Comment)
	}
}

func TestTokenizeSimplePair(t *testing.T) {
	toks, err := Tokenize([]string{"FOO=bar"})
	if err != nil {
		t.Fatal(err)
	}
	if toks[0].Key != "FOO" || toks[0].Value != "bar" {
		t.Errorf("unexpected token: %+v", toks[0])
	}
}

func TestTokenizeQuotedValue(t *testing.T) {
	toks, err := Tokenize([]string{`DB_URL="postgres://localhost/mydb"`})
	if err != nil {
		t.Fatal(err)
	}
	if toks[0].Value != "postgres://localhost/mydb" {
		t.Errorf("expected unquoted value, got %q", toks[0].Value)
	}
}

func TestTokenizeInlineComment(t *testing.T) {
	toks, err := Tokenize([]string{"PORT=8080 # http port"})
	if err != nil {
		t.Fatal(err)
	}
	if toks[0].Value != "8080" {
		t.Errorf("expected value '8080', got %q", toks[0].Value)
	}
	if toks[0].Comment != "http port" {
		t.Errorf("expected inline comment 'http port', got %q", toks[0].Comment)
	}
}

func TestTokenizeHashInsideQuotes(t *testing.T) {
	toks, err := Tokenize([]string{`MSG="hello # world"`})
	if err != nil {
		t.Fatal(err)
	}
	if toks[0].Value != "hello # world" {
		t.Errorf("hash inside quotes should not split, got %q", toks[0].Value)
	}
	if toks[0].Comment != "" {
		t.Errorf("expected no inline comment, got %q", toks[0].Comment)
	}
}

func TestTokenizeMissingEquals(t *testing.T) {
	_, err := Tokenize([]string{"BADLINE"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestTokenizeString(t *testing.T) {
	src := "KEY1=val1\n# comment\n\nKEY2=val2"
	toks, err := TokenizeString(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) != 4 {
		t.Fatalf("expected 4 tokens, got %d", len(toks))
	}
	if toks[0].Type != TokenPair || toks[1].Type != TokenComment ||
		toks[2].Type != TokenBlank || toks[3].Type != TokenPair {
		t.Errorf("unexpected token types: %v %v %v %v",
			toks[0].Type, toks[1].Type, toks[2].Type, toks[3].Type)
	}
}

package envfile

import (
	"fmt"
	"strings"
)

// TokenType represents the kind of token parsed from an env line.
type TokenType int

const (
	TokenComment TokenType = iota
	TokenBlank
	TokenPair
)

// Token represents a single lexical unit from an env file line.
type Token struct {
	Type    TokenType
	Key     string
	Value   string
	Comment string
	Raw     string
}

// Tokenize parses a slice of raw lines into Tokens.
// It handles comments, blank lines, quoted values, and inline comments.
func Tokenize(lines []string) ([]Token, error) {
	tokens := make([]Token, 0, len(lines))
	for i, line := range lines {
		tok, err := tokenizeLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		tokens = append(tokens, tok)
	}
	return tokens, nil
}

// TokenizeString tokenizes a multi-line env string.
func TokenizeString(src string) ([]Token, error) {
	lines := strings.Split(src, "\n")
	return Tokenize(lines)
}

func tokenizeLine(raw string) (Token, error) {
	trimmed := strings.TrimSpace(raw)

	if trimmed == "" {
		return Token{Type: TokenBlank, Raw: raw}, nil
	}
	if strings.HasPrefix(trimmed, "#") {
		return Token{Type: TokenComment, Comment: strings.TrimPrefix(trimmed, "#"), Raw: raw}, nil
	}

	eqIdx := strings.IndexByte(trimmed, '=')
	if eqIdx < 0 {
		return Token{}, fmt.Errorf("missing '=' in %q", raw)
	}

	key := strings.TrimSpace(trimmed[:eqIdx])
	rest := trimmed[eqIdx+1:]

	value, inline := splitInlineComment(rest)

	return Token{
		Type:    TokenPair,
		Key:     key,
		Value:   unquote(strings.TrimSpace(value)),
		Comment: inline,
		Raw:     raw,
	}, nil
}

// splitInlineComment separates value from an optional trailing # comment.
// Quoted sections are respected so that # inside quotes is not treated as a comment.
func splitInlineComment(s string) (value, comment string) {
	inQuote := false
	var quoteChar byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inQuote {
			if c == quoteChar {
				inQuote = false
			}
			continue
		}
		if c == '"' || c == '\'' {
			inQuote = true
			quoteChar = c
			continue
		}
		if c == '#' {
			return s[:i], strings.TrimSpace(s[i+1:])
		}
	}
	return s, ""
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

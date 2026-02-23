package lexer

import (
	"testing"
)

func TestFastLexer_BackwardCompatibility(t *testing.T) {
	source := `
	dialect:english;
	// This is a comment
	func main() {
		print("Hello");
	}
	`

	// Default behavior: No whitespace, no comments
	lex, _, err := NewFastLexer(source)
	if err != nil {
		t.Fatalf("Failed to create lexer: %v", err)
	}

	tokens, err := lex.GetAllTokens()
	if err != nil {
		t.Fatalf("Failed to tokenize: %v", err)
	}

	for _, token := range tokens {
		if token.Kind == Whitespace || token.Kind == Comment || token.Kind == Newline {
			t.Errorf("Unexpected token kind %s found in default mode", token.Kind)
		}
	}
}

func TestFastLexer_WithWhitespaceAndComments(t *testing.T) {
	source := `dialect:english; // comment`

	lex, _, err := NewFastLexer(source)
	if err != nil {
		t.Fatalf("Failed to create lexer: %v", err)
	}

	// Enable flags
	lex.ShowWhitespace = true
	lex.ShowComments = true

	tokens, err := lex.GetAllTokens()
	if err != nil {
		t.Fatalf("Failed to tokenize: %v", err)
	}

	hasWhitespace := false
	hasComment := false

	for _, token := range tokens {
		if token.Kind == Whitespace {
			hasWhitespace = true
		}
		if token.Kind == Comment {
			hasComment = true
		}
	}

	if !hasWhitespace {
		t.Error("Expected Whitespace token, but none found")
	}
	if !hasComment {
		t.Error("Expected Comment token, but none found")
	}
}

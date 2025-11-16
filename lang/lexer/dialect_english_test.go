package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservedEnglish(t *testing.T) {
	expected := map[string]TokenKind{
		"if":     KeywordIf,
		"else":   KeywordElse,
		"struct": Struct,
		"let":    Let,
		"const":  Const,
		"int":    TypeInt,
		"float":  TypeFloat,
		"string": TypeString,
	}
	english := English{}

	assert.Equal(t, expected, english.Reserved())
}

func TestIntegerTokenization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TokenKind
		value    string
	}{
		{"simple integer", "dialect:english; 42", Integer, "42"},
		{"zero", "dialect:english; 0", Integer, "0"},
		{"large integer", "dialect:english; 123456789", Integer, "123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := Tokenize(tt.input)
			// tokens: [dialect, :, english, ;, NUMBER, EOF]
			assert.True(t, len(tokens) >= 5, "Expected at least 5 tokens")
			assert.Equal(t, tt.expected, tokens[4].Kind)
			assert.Equal(t, tt.value, tokens[4].Value)
		})
	}
}

func TestFloatTokenization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TokenKind
		value    string
	}{
		{"simple float", "dialect:english; 3.14", Float, "3.14"},
		{"float with zero", "dialect:english; 0.5", Float, "0.5"},
		{"large float", "dialect:english; 123.456", Float, "123.456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := Tokenize(tt.input)
			// tokens: [dialect, :, english, ;, NUMBER, EOF]
			assert.True(t, len(tokens) >= 5, "Expected at least 5 tokens")
			assert.Equal(t, tt.expected, tokens[4].Kind)
			assert.Equal(t, tt.value, tokens[4].Value)
		})
	}
}

func TestIntegerVsFloatTokenization(t *testing.T) {
	input := "dialect:english; let x = 42; let y = 3.14; let z = 100; let w = 2.5;"
	tokens, _ := Tokenize(input)

	// Find the numeric tokens
	var numbers []Token
	for _, tok := range tokens {
		if tok.Kind == Integer || tok.Kind == Float {
			numbers = append(numbers, tok)
		}
	}

	assert.Equal(t, 4, len(numbers), "Expected 4 numeric tokens")
	assert.Equal(t, Integer, numbers[0].Kind)
	assert.Equal(t, "42", numbers[0].Value)
	assert.Equal(t, Float, numbers[1].Kind)
	assert.Equal(t, "3.14", numbers[1].Value)
	assert.Equal(t, Integer, numbers[2].Kind)
	assert.Equal(t, "100", numbers[2].Value)
	assert.Equal(t, Float, numbers[3].Kind)
	assert.Equal(t, "2.5", numbers[3].Value)
}

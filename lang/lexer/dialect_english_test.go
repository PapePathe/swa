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
			// Skip dialect tokens, semicolon - the number should be at index 3
			assert.True(t, len(tokens) >= 4, "Expected at least 4 tokens")
			assert.Equal(t, tt.expected, tokens[3].Kind)
			assert.Equal(t, tt.value, tokens[3].Value)
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
			// Skip dialect tokens, semicolon - the number should be at index 3
			assert.True(t, len(tokens) >= 4, "Expected at least 4 tokens")
			assert.Equal(t, tt.expected, tokens[3].Kind)
			assert.Equal(t, tt.value, tokens[3].Value)
		})
	}
}

func TestIntegerVsFloatTokenization(t *testing.T) {
	input := "dialect:english; 42 3.14 100 2.5"
	tokens, _ := Tokenize(input)

	// Skip dialect:english; tokens (3 tokens), then we have 4 numbers + EOF
	assert.True(t, len(tokens) >= 8, "Expected at least 8 tokens")
	assert.Equal(t, Integer, tokens[3].Kind)
	assert.Equal(t, "42", tokens[3].Value)
	assert.Equal(t, Float, tokens[4].Kind)
	assert.Equal(t, "3.14", tokens[4].Value)
	assert.Equal(t, Integer, tokens[5].Kind)
	assert.Equal(t, "100", tokens[5].Value)
	assert.Equal(t, Float, tokens[6].Kind)
	assert.Equal(t, "2.5", tokens[6].Value)
}

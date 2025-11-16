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
		{"simple integer", "42", Integer, "42"},
		{"zero", "0", Integer, "0"},
		{"large integer", "123456789", Integer, "123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := &Lexer{
				Tokens:        make([]Token, 0),
				patterns:      English{}.Patterns(),
				reservedWords: English{}.Reserved(),
				source:        tt.input,
			}
			tokens := lex.Tokenize()
			assert.Equal(t, 2, len(tokens), "Expected 2 tokens (value + EOF)")
			assert.Equal(t, tt.expected, tokens[0].Kind)
			assert.Equal(t, tt.value, tokens[0].Value)
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
		{"simple float", "3.14", Float, "3.14"},
		{"float with zero", "0.5", Float, "0.5"},
		{"large float", "123.456", Float, "123.456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := &Lexer{
				Tokens:        make([]Token, 0),
				patterns:      English{}.Patterns(),
				reservedWords: English{}.Reserved(),
				source:        tt.input,
			}
			tokens := lex.Tokenize()
			assert.Equal(t, 2, len(tokens), "Expected 2 tokens (value + EOF)")
			assert.Equal(t, tt.expected, tokens[0].Kind)
			assert.Equal(t, tt.value, tokens[0].Value)
		})
	}
}

func TestIntegerVsFloatTokenization(t *testing.T) {
	input := "42 3.14 100 2.5"
	lex := &Lexer{
		Tokens:        make([]Token, 0),
		patterns:      English{}.Patterns(),
		reservedWords: English{}.Reserved(),
		source:        input,
	}
	tokens := lex.Tokenize()

	assert.Equal(t, 5, len(tokens), "Expected 5 tokens (4 numbers + EOF)")
	assert.Equal(t, Integer, tokens[0].Kind)
	assert.Equal(t, "42", tokens[0].Value)
	assert.Equal(t, Float, tokens[1].Kind)
	assert.Equal(t, "3.14", tokens[1].Value)
	assert.Equal(t, Integer, tokens[2].Kind)
	assert.Equal(t, "100", tokens[2].Value)
	assert.Equal(t, Float, tokens[3].Kind)
	assert.Equal(t, "2.5", tokens[3].Value)
}

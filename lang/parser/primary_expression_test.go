package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func TestParseIntegerExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"simple integer", "dialect:english; 42", 42},
		{"zero", "dialect:english; 0", 0},
		{"large integer", "dialect:english; 123456", 123456},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := lexer.Tokenize(tt.input)
			parser := New(tokens)

			// Skip to the number token (after dialect:english;)
			parser.advance() // dialect
			parser.advance() // :
			parser.advance() // english
			parser.advance() // ;

			expr, err := ParsePrimaryExpression(parser)
			assert.NoError(t, err)

			intExpr, ok := expr.(ast.IntegerExpression)
			assert.True(t, ok, "Expected IntegerExpression")
			assert.Equal(t, tt.expected, intExpr.Value)
		})
	}
}

func TestParseFloatExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"simple float", "dialect:english; 3.14", 3.14},
		{"float with zero", "dialect:english; 0.5", 0.5},
		{"large float", "dialect:english; 123.456", 123.456},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := lexer.Tokenize(tt.input)
			parser := New(tokens)

			// Skip to the number token (after dialect:english;)
			parser.advance() // dialect
			parser.advance() // :
			parser.advance() // english
			parser.advance() // ;

			expr, err := ParsePrimaryExpression(parser)
			assert.NoError(t, err)

			floatExpr, ok := expr.(ast.FloatExpression)
			assert.True(t, ok, "Expected FloatExpression")
			assert.Equal(t, tt.expected, floatExpr.Value)
		})
	}
}

func TestParseIntegerType(t *testing.T) {
	input := "dialect:english; int"
	tokens, _ := lexer.Tokenize(input)
	parser := New(tokens)

	// Skip to type token
	parser.advance() // dialect
	parser.advance() // :
	parser.advance() // english
	parser.advance() // ;

	typ, _ := parseType(parser, DefaultBindingPower)

	intType, ok := typ.(ast.IntegerType)
	assert.True(t, ok, "Expected IntegerType")
	assert.Equal(t, ast.DataTypeInteger, intType.Value())
}

func TestParseFloatType(t *testing.T) {
	input := "dialect:english; float"
	tokens, _ := lexer.Tokenize(input)
	parser := New(tokens)

	// Skip to type token
	parser.advance() // dialect
	parser.advance() // :
	parser.advance() // english
	parser.advance() // ;

	typ, _ := parseType(parser, DefaultBindingPower)

	floatType, ok := typ.(ast.FloatType)
	assert.True(t, ok, "Expected FloatType")
	assert.Equal(t, ast.DataTypeFloat, floatType.Value())
}

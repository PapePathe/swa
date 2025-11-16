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
		{"simple integer", "42", 42},
		{"zero", "0", 0},
		{"large integer", "123456", 123456},
		{"negative integer", "-5", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.New(tt.input, lexer.English{})
			tokens := lex.Tokenize()
			parser := New(tokens)

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
		{"simple float", "3.14", 3.14},
		{"float with zero", "0.5", 0.5},
		{"large float", "123.456", 123.456},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.New(tt.input, lexer.English{})
			tokens := lex.Tokenize()
			parser := New(tokens)

			expr, err := ParsePrimaryExpression(parser)
			assert.NoError(t, err)

			floatExpr, ok := expr.(ast.FloatExpression)
			assert.True(t, ok, "Expected FloatExpression")
			assert.Equal(t, tt.expected, floatExpr.Value)
		})
	}
}

func TestParseIntegerType(t *testing.T) {
	input := "int"
	lex := lexer.New(input, lexer.English{})
	tokens := lex.Tokenize()
	parser := New(tokens)

	typ, _ := parseType(parser, DefaultBindingPower)

	intType, ok := typ.(ast.IntegerType)
	assert.True(t, ok, "Expected IntegerType")
	assert.Equal(t, ast.DataTypeInteger, intType.Value())
}

func TestParseFloatType(t *testing.T) {
	input := "float"
	lex := lexer.New(input, lexer.English{})
	tokens := lex.Tokenize()
	parser := New(tokens)

	typ, _ := parseType(parser, DefaultBindingPower)

	floatType, ok := typ.(ast.FloatType)
	assert.True(t, ok, "Expected FloatType")
	assert.Equal(t, ast.DataTypeFloat, floatType.Value())
}

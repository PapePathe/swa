package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservedEnglish(t *testing.T) {
	expected := map[string]TokenKind{
		"const":    Const,
		"dialect":  DialectDeclaration,
		"else":     KeywordElse,
		"error":    TypeError,
		"float":    TypeFloat,
		"func":     Function,
		"if":       KeywordIf,
		"int":      TypeInt,
		"int64":    TypeInt64,
		"let":      Let,
		"print":    Print,
		"return":   Return,
		"string":   TypeString,
		"start":    Main,
		"struct":   Struct,
		"variadic": Variadic,
		"while":    KeywordWhile,
		"zero":     Zero,
	}
	english := English{}

	assert.Equal(t, expected, english.Reserved())
}

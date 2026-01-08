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
	}
	english := English{}

	assert.Equal(t, expected, english.Reserved())
}

func TestPatternsEnglish(t *testing.T) {
	english := English{}

	assert.Equal(t, 48, len(english.Patterns()))
}

package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservedEnglish(t *testing.T) {
	expected := map[string]TokenKind{
		"const":    Const,
		"dialect":  DialectDeclaration,
		"func":     Function,
		"else":     KeywordElse,
		"if":       KeywordIf,
		"while":    KeywordWhile,
		"let":      Let,
		"start":    Main,
		"print":    Print,
		"return":   Return,
		"struct":   Struct,
		"true":     True,
		"false":    False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"int":      TypeInt,
		"int64":    TypeInt64,
		"string":   TypeString,
		"error":    TypeError,
		"variadic": Variadic,
		"zero":     Zero,
		"make":     Make,
		"count":    Len,
		"cap":      Cap,
		"append":   Append,
	}

	english := English{}

	assert.Equal(t, expected, english.Reserved())
}

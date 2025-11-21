package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservedEnglish(t *testing.T) {
	expected := map[string]TokenKind{
		"if":     KeywordIf,
		"else":   KeywordElse,
		"func":   Function,
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

func TestPatternsEnglish(t *testing.T) {
	english := English{}

	assert.Equal(t, 46, len(english.Patterns()))

}

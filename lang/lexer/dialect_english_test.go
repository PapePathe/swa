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
	}
	english := English{}

	assert.Equal(t, expected, english.Reserved())
}

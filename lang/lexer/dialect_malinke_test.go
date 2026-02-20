package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservedMalinke(t *testing.T) {
	expected := map[string]TokenKind{
		"tintin":    Const,
		"kan":       DialectDeclaration,
		"baara":     Function,
		"wala":      KeywordElse,
		"ni":        KeywordIf,
		"tuma":      KeywordWhile,
		"atö":       Let,
		"daminen":   Main,
		"yira":      Print,
		"segin":     Return,
		"joyoro":    Struct,
		"jatelen":   TypeFloat,
		"jate":      TypeInt,
		"jate64":    TypeInt64,
		"seben":     TypeString,
		"fili":      TypeError,
		"caaman":    Variadic,
		"foy":       Zero,
		"tinye":     True,
		"wouya":     False,
		"tinyejate": TypeBool,
	}

	malinke := Malinke{}
	reserved := malinke.Reserved()

	assert.Equal(t, expected, reserved)
}

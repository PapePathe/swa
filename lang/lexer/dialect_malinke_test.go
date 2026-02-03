package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservedMalinke(t *testing.T) {
	expected := map[string]TokenKind{
		"ni":     KeywordIf,
		"nii":    KeywordElse,
		"struct": Struct,
		"let":    Let,
		"const":  Const,
		"fèndo":  TypeInt,
		"erreur": TypeError,
	}

	malinke := Malinke{}
	reserved := malinke.Reserved()

	assert.Equal(t, expected, reserved)
}

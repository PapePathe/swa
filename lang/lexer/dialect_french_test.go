package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservedFrench(t *testing.T) {
	expected := map[string]TokenKind{
		"si":        KeywordIf,
		"sinon":     KeywordElse,
		"structure": Struct,
		"fonction":  Function,
		"variable":  Let,
		"constante": Const,
		"entier":    TypeInt,
		"decimal":   TypeFloat,
		"chaine":    TypeString,
	}

	French := French{}
	reserved := French.Reserved()

	assert.Equal(t, expected, reserved)
}

func TestPatternsFrench(t *testing.T) {
	french := French{}

	assert.Equal(t, 46, len(french.Patterns()))
}

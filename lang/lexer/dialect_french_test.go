package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservedFrench(t *testing.T) {
	expected := map[string]TokenKind{
		"afficher":   Print,
		"chaine":     TypeString,
		"constante":  Const,
		"decimal":    TypeFloat,
		"demarrer":   Main,
		"dialecte":   DialectDeclaration,
		"entier":     TypeInt,
		"entier64":   TypeInt64,
		"fonction":   Function,
		"retourner":  Return,
		"si":         KeywordIf,
		"sinon":      KeywordElse,
		"structure":  Struct,
		"tantque":    KeywordWhile,
		"variable":   Let,
		"variadique": Variadic,
	}

	French := French{}
	reserved := French.Reserved()

	assert.Equal(t, expected, reserved)
}

func TestPatternsFrench(t *testing.T) {
	french := French{}

	assert.Equal(t, 48, len(french.Patterns()))
}

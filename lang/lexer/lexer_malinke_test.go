package lexer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTokenizeMalinke ...
func TestTokenizeMalinke(t *testing.T) {
	bytes, err := os.ReadFile("../examples/malinke/age_calculator.swa")
	require.NoError(t, err)

	source := string(bytes)
	expecedTokens := []Token{
		{Value: "fèndo", Kind: TypeInt},
		{Value: "x", Kind: Identifier},
		{Value: "=", Kind: Assignment},
		{Value: "10", Kind: Number},
		{Value: ";", Kind: SemiColon},
		{Value: "ni", Kind: KeywordIf},
		{Value: "(", Kind: OpenParen},
		{Value: "x", Kind: Identifier},
		{Value: ">", Kind: GreaterThan},
		{Value: "10", Kind: Number},
		{Value: ")", Kind: CloseParen},
		{Value: "{", Kind: OpenCurly},
		{Value: "afo", Kind: Identifier},
		{Value: "(", Kind: OpenParen},
		{Value: "\"Isoma\"", Kind: String},
		{Value: ")", Kind: CloseParen},
		{Value: ";", Kind: SemiColon},
		{Value: "}", Kind: CloseCurly},
		{Value: "nii", Kind: KeywordElse},
		{Value: "{", Kind: OpenCurly},
		{Value: "afo", Kind: Identifier},
		{Value: "(", Kind: OpenParen},
		{Value: "\"Inoura\"", Kind: String},
		{Value: ")", Kind: CloseParen},
		{Value: ";", Kind: SemiColon},
		{Value: "}", Kind: CloseCurly},
		{Value: "EOF", Kind: EOF},
	}

	assert.Equal(t, Tokenize(source), expecedTokens)
}

package lexer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizeMalinke(t *testing.T) {
	bytes, err := os.ReadFile("../examples/malinke.swa")
	require.NoError(t, err)

	source := string(bytes)

	assert.Equal(
		t,
		Tokenize(source),
		[]Token{
			{Value: "fèndo", Kind: TYPE_INT},
			{Value: "x", Kind: IDENTIFIER},
			{Value: "=", Kind: ASSIGNMENT},
			{Value: "10", Kind: NUMBER},
			{Value: ";", Kind: SEMI_COLON},
			{Value: "ni", Kind: KEYWORD_IF},
			{Value: "(", Kind: OPEN_PAREN},
			{Value: "x", Kind: IDENTIFIER},
			{Value: ">", Kind: GREATER_THAN},
			{Value: "10", Kind: NUMBER},
			{Value: ")", Kind: CLOSE_PAREN},
			{Value: "{", Kind: OPEN_CURLY},
			{Value: "afo", Kind: IDENTIFIER},
			{Value: "(", Kind: OPEN_PAREN},
			{Value: "\"Isoma\"", Kind: STRING},
			{Value: ")", Kind: CLOSE_PAREN},
			{Value: ";", Kind: SEMI_COLON},
			{Value: "}", Kind: CLOSE_CURLY},
			{Value: "nii", Kind: KEYWORD_ELSE},
			{Value: "{", Kind: OPEN_CURLY},
			{Value: "afo", Kind: IDENTIFIER},
			{Value: "(", Kind: OPEN_PAREN},
			{Value: "\"Inoura\"", Kind: STRING},
			{Value: ")", Kind: CLOSE_PAREN},
			{Value: ";", Kind: SEMI_COLON},
			{Value: "}", Kind: CLOSE_CURLY},
			{Value: "EOF", Kind: EOF},
		},
	)
}

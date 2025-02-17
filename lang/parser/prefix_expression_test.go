package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrefixExpression(t *testing.T) {
	st := Parse(lexer.Tokenize("-44;"))

	assert.Equal(t, st, ast.PrefixExpression{})
}

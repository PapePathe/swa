package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrintStatement(t *testing.T) {
	expectedSimpleAst := ast.BlockStatement{
		Body: []ast.Statement{
			ast.PrintStatetement{
				Values: []ast.Expression{
					ast.StringExpression{Value: "une chaine de caractères"},
				},
			},
		},
	}
	expectedNumberAst := ast.BlockStatement{
		Body: []ast.Statement{
			ast.PrintStatetement{
				Values: []ast.Expression{
					ast.NumberExpression{Value: 10.0000},
				},
			},
		},
	}
	expectedIdentifierAst := ast.BlockStatement{
		Body: []ast.Statement{
			ast.PrintStatetement{
				Values: []ast.Expression{
					ast.SymbolExpression{Value: "x"},
				},
			},
		},
	}

	t.Run("french", func(t *testing.T) {
		t.Run("simple", func(t *testing.T) {
			source := `dialect:french; afficher("une chaine de caractères");`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedSimpleAst, result)
		})
		t.Run("with a number", func(t *testing.T) {
			source := `dialect:french; afficher(10);`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedNumberAst, result)
		})
		t.Run("with an identifier", func(t *testing.T) {
			source := `dialect:french; afficher(x);`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedIdentifierAst, result)
		})
	})
	t.Run("english", func(t *testing.T) {
		t.Run("simple", func(t *testing.T) {
			source := `dialect:english; print("une chaine de caractères");`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedSimpleAst, result)
		})
		t.Run("with a number", func(t *testing.T) {
			source := `dialect:english; print(10);`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedNumberAst, result)
		})
		t.Run("with an identifier", func(t *testing.T) {
			source := `dialect:english; print(x);`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedIdentifierAst, result)
		})
	})
	t.Run("malinke", func(t *testing.T) {
		t.Run("simple", func(t *testing.T) {
			source := `dialect:malinke; afo("une chaine de caractères");`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedSimpleAst, result)
		})
		t.Run("with a number", func(t *testing.T) {
			source := `dialect:malinke; afo(10);`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedNumberAst, result)
		})
		t.Run("with an identifier", func(t *testing.T) {
			source := `dialect:malinke; afo(x);`
			result := parseSourceCode(t, source)

			assert.Equal(t, expectedIdentifierAst, result)
		})
	})
}

func parseSourceCode(t *testing.T, src string) ast.BlockStatement {
	t.Helper()

	return Parse(lexer.Tokenize(src))
}

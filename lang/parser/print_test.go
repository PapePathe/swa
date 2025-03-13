package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/values"
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
	})
}

func parseSourceCode(t *testing.T, src string) ast.BlockStatement {
	t.Helper()

	return Parse(lexer.Tokenize(src))
}

func TestEvaluatePrintStatement(t *testing.T) {
	t.Run("with static strings", func(t *testing.T) {
		stmt := ast.PrintStatetement{
			Values: []ast.Expression{ast.StringExpression{Value: "Hello"}},
		}
		err, val := stmt.Evaluate(ast.NewScope(nil))

		assert.NoError(t, err)
		assert.Equal(t, val, values.StringValue{Value: "Hello"})
	})

	t.Run("with variables", func(t *testing.T) {
		scope := ast.NewScope(nil)
		scope.Set("x", values.StringValue{Value: "HI"})

		stmt := ast.PrintStatetement{
			Values: []ast.Expression{ast.SymbolExpression{Value: "x"}},
		}
		err, val := stmt.Evaluate(scope)

		assert.NoError(t, err)
		assert.Equal(t, val, values.StringValue{Value: "HI"})
	})

	t.Run("with numbers", func(t *testing.T) {
		stmt := ast.PrintStatetement{
			Values: []ast.Expression{ast.NumberExpression{Value: 10}},
		}
		err, val := stmt.Evaluate(ast.NewScope(nil))

		assert.NoError(t, err)
		assert.Equal(t, val, values.StringValue{Value: "10.000000"})
	})
}

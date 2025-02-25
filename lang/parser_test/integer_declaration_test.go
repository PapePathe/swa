package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func TestImplicitIntegerDeclaration(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      // dialect=malinke;
      let nombre = -44.5;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:       "nombre",
				IsConstant: false,
				Value: ast.PrefixExpression{
					Operator:        lexer.Token{Value: "-", Kind: 22},
					RightExpression: ast.NumberExpression{Value: 44.5},
				},
				ExplicitType: ast.Type(nil), // exlicit type should be decimal
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestExplicitIntegerDeclaration(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      // dialect=malinke;
			let nombre : int = 4 +3;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:       "nombre",
				IsConstant: false,
				Value: ast.BinaryExpression{
					Left:     ast.NumberExpression{Value: 4},
					Right:    ast.NumberExpression{Value: 3},
					Operator: lexer.Token{Value: "+", Kind: 30},
				},
				ExplicitType: ast.SymbolType{Name: "int"},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestIntegerDeclarationWithComplexExpression(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      // dialect=malinke;
			let resultat : int = nombre * (45 - -5);
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:       "resultat",
				IsConstant: false,
				Value: ast.BinaryExpression{
					Left: ast.SymbolExpression{Value: "nombre"},
					Right: ast.BinaryExpression{
						Left: ast.NumberExpression{Value: 45},
						Right: ast.PrefixExpression{
							Operator:        lexer.Token{Value: "-", Kind: 22},
							RightExpression: ast.NumberExpression{Value: 5},
						},
						Operator: lexer.Token{Value: "-", Kind: 22},
					},
					Operator: lexer.Token{Value: "*", Kind: 32},
				},
				ExplicitType: ast.SymbolType{Name: "int"},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestConstantIntegerExplicitDeclaration(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      // dialect=malinke;
			const resultat : int = 20;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:         "resultat",
				IsConstant:   true,
				Value:        ast.NumberExpression{Value: 20},
				ExplicitType: ast.SymbolType{Name: "int"},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestConstantIntegerImplicitDeclaration(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      // dialect=malinke;
			const resultat  = 20;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:         "resultat",
				IsConstant:   true,
				Value:        ast.NumberExpression{Value: 20},
				ExplicitType: nil,
			},
		},
	}

	assert.Equal(t, result, expected)
}

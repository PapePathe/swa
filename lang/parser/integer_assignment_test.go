package parser_test

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerAssigment(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
		  dialect:malinke;
      nombre = 4;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "=", Kind: 9, Name: "ASSIGNMENT"},
					Assignee: ast.SymbolExpression{Value: "nombre"},
					Value:    ast.NumberExpression{Value: 4},
				},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithInfixOperator(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  		dialect:malinke;
      nombre = -4;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "=", Kind: 9, Name: "ASSIGNMENT"},
					Assignee: ast.SymbolExpression{Value: "nombre"},
					Value: ast.PrefixExpression{
						Operator:        lexer.Token{Value: "-", Kind: 22, Name: "MINUS"},
						RightExpression: ast.NumberExpression{Value: 4},
					},
				},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithAdditionToSelf(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  		dialect:malinke;
      nombre += 4;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "+=", Kind: 31, Name: "PLUS_EQUAL"},
					Assignee: ast.SymbolExpression{Value: "nombre"},
					Value:    ast.NumberExpression{Value: 4},
				},
			},
		},
	}

	assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithMultiplicationToSelf(t *testing.T) {
	// result := parser.Parse(lexer.Tokenize(`
	//     // dialect=malinke;
	//     nombre *= 4;
	// `))

	// expected := ast.BlockStatement{
	// 	Body: []ast.Statement{
	// 		ast.ExpressionStatement{
	// 			Exp: ast.AssignmentExpression{
	// 				Operator: lexer.Token{Value: "+=", Kind: 31},
	// 				Assignee: ast.SymbolExpression{Value: "nombre"},
	// 				Value:    ast.NumberExpression{Value: 4},
	// 			},
	// 		},
	// 	},
	// }

	// assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithSubstractionToSelf(t *testing.T) {
	// result := parser.Parse(lexer.Tokenize(`
	//     // dialect=malinke;
	//     nombre *= 4;
	// `))

	// expected := ast.BlockStatement{
	// 	Body: []ast.Statement{
	// 		ast.ExpressionStatement{
	// 			Exp: ast.AssignmentExpression{
	// 				Operator: lexer.Token{Value: "+=", Kind: 31},
	// 				Assignee: ast.SymbolExpression{Value: "nombre"},
	// 				Value:    ast.NumberExpression{Value: 4},
	// 			},
	// 		},
	// 	},
	// }

	// assert.Equal(t, result, expected)
}

func TestIntegerAssigmentWithExpression(t *testing.T) {
	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:       "multiply",
				IsConstant: true,
				Value: ast.BinaryExpression{
					Left: ast.NumberExpression{Value: 45.2},
					Right: ast.BinaryExpression{
						Left:     ast.NumberExpression{Value: 5},
						Right:    ast.NumberExpression{Value: 4},
						Operator: lexer.Token{Value: "*", Kind: 32, Name: "STAR"},
					},
					Operator: lexer.Token{Value: "+", Kind: 30, Name: "PLUS"},
				},
				ExplicitType: ast.Type(nil),
			},
		},
	}

	result := parser.Parse(lexer.Tokenize(`
			 dialect:malinke;
       const multiply = 45.2 + 5 * 4;
	 `))

	assert.Equal(t, result, expected)
}

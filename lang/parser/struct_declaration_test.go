package parser_test

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedAstForTestStructDeclaration = ast.BlockStatement{
	Body: []ast.Statement{
		ast.StructDeclarationStatement{
			Name: "Rectangle",
			Properties: map[string]ast.StructProperty{
				"height": {PropType: ast.SymbolType{Name: "float"}},
				"name":   {PropType: ast.SymbolType{Name: "string"}},
				"width":  {PropType: ast.SymbolType{Name: "number"}},
			},
		},
	},
}

func TestStructDeclarationMalinke(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:malinke;
      struct Rectangle {
        width: number,
        height: float,
        name: string,
      }
	`))
	assert.Equal(t, result, expectedAstForTestStructDeclaration)
}

func TestStructDeclarationEnglish(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:english;
      struct Rectangle {
        width: number,
        height: float,
        name: string,
      }
	`))
	assert.Equal(t, result, expectedAstForTestStructDeclaration)
}

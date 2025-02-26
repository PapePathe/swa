package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func TestStructDeclaration(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:malinke;
      struct Rectangle {
        width: number,
        height: float,
        name: string,
      }
	`))

	expected := ast.BlockStatement{
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

	assert.Equal(t, result, expected)
}

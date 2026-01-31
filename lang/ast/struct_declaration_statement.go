package ast

import (
	"fmt"
	"swahili/lang/lexer"
)

type StructDeclarationStatement struct {
	Name       string
	Properties []string
	Types      []Type
	Tokens     []lexer.Token
	SwaType    Type
}

var _ Statement = (*StructDeclarationStatement)(nil)

func (stmt StructDeclarationStatement) PropertyIndex(name string) (error, int) {
	for propIndex, propName := range stmt.Properties {
		if propName == name {
			return nil, propIndex
		}
	}

	return fmt.Errorf("Property with name (%s) does not exist on struct %s", name, stmt.Name), 0
}

func (stmt *StructDeclarationStatement) Accept(g CodeGenerator) error {
	return g.VisitStructDeclaration(stmt)
}

func (stmt StructDeclarationStatement) TokenStream() []lexer.Token {
	return stmt.Tokens
}

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
}

var _ Statement = (*StructDeclarationStatement)(nil)

func (sd StructDeclarationStatement) PropertyIndex(name string) (error, int) {
	for propIndex, propName := range sd.Properties {
		if propName == name {
			return nil, propIndex
		}
	}

	return fmt.Errorf("Property with name (%s) does not exist on struct %s", name, sd.Name), 0
}

func (sd StructDeclarationStatement) Accept(g CodeGenerator) error {
	return g.VisitStructDeclaration(&sd)
}

func (expr StructDeclarationStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

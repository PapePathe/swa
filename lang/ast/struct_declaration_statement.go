package ast

import (
	"swahili/lang/values"
)

type StructProperty struct {
	PropType Type
}
type StructDeclarationStatement struct {
	Name       string
	Properties map[string]StructProperty
}

var _ Statement = (*StructDeclarationStatement)(nil)

func (cs StructDeclarationStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (s StructDeclarationStatement) statement() {}

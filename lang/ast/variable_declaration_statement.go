package ast

import "swahili/lang/values"

// VarDeclarationStatement ...
type VarDeclarationStatement struct {
	// The name of the variable
	Name string
	// Wether or not the variable is a constant
	IsConstant bool
	// The value assigned to the variable
	Value Expression
	// The explicit type of the variable
	ExplicitType Type
}

var _ Statement = (*VarDeclarationStatement)(nil)

func (v VarDeclarationStatement) Evaluate(s *Scope) (error, values.Value) {
	_, val := v.Value.Evaluate(s)

	s.Set(v.Name, val)

	return nil, nil
}

func (bs VarDeclarationStatement) statement() {}

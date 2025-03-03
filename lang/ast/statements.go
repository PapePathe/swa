package ast

import (
	"fmt"
	"swahili/lang/values"
)

type ConditionalStatetement struct {
	Condition Expression
	Success   BlockStatement
	Failure   BlockStatement
}

var _ Statement = (*ConditionalStatetement)(nil)

func (cs ConditionalStatetement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
func (cs ConditionalStatetement) statement() {}

// BlockStatement ...
type BlockStatement struct {
	// The body of the block statement
	Body []Statement
}

var _ Statement = (*BlockStatement)(nil)

func (cs BlockStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (bs BlockStatement) statement() {}

// ExpressionStatement ...
type ExpressionStatement struct {
	// The expression
	Exp Expression
}

var _ Statement = (*ExpressionStatement)(nil)

func (cs ExpressionStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (bs ExpressionStatement) statement() {}

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
	fmt.Println("Evaluating variable declaration statement")

	_, val := v.Value.Evaluate(s)

	s.Set(v.Name, val)

	return nil, nil
}

func (bs VarDeclarationStatement) statement() {}

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

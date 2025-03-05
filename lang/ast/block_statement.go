package ast

import "swahili/lang/values"

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

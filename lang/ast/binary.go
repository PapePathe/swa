package ast

import (
	"fmt"
	"swahili/lang/lexer"
	"swahili/lang/values"
)

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

var _ Expression = (*BinaryExpression)(nil)

func (BinaryExpression) expression() {}

func (be BinaryExpression) Evaluate(s *Scope) (error, values.Value) {
	err, left := be.Left.Evaluate(s)
	if err != nil {
		return fmt.Errorf("Errored while evaluating left : <%s>", err), nil
	}

	err, right := be.Right.Evaluate(s)
	if err != nil {
		return fmt.Errorf("Errored while evaluating right : <%s>", err), nil
	}

	switch be.Operator.Kind {
	case lexer.Star:
		return nil, nil
	case lexer.Minus:
		return nil, nil

	default:
		panic(fmt.Sprintf("unknown operator: %s", be.Operator))
	}
}

package ast

import "swahili/lang/values"

// Statement ...
type Statement interface {
	statement()
	Evaluate(s *Scope) (error, values.Value)
}

// Expression ...
type Expression interface {
	expression()
	Evaluate(s *Scope) (error, values.Value)
}

type Type interface {
	_type()
}

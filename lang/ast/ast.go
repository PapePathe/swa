package ast

// Statement ...
type Statement interface {
	statement()
}

// Expression ...
type Expression interface {
	expression()
}

type Type interface {
	_type()
}

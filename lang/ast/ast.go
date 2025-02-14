package ast

type Statement interface {
	statement()
}

type Expression interface {
	expression()
}

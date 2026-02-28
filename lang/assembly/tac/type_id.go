package tac

import "swahili/lang/ast"

type TypeID struct {
	T ast.Type
}

var _ InstArg = (*TypeID)(nil)

func (t *TypeID) InstructionArg() string {
	return t.T.String()
}

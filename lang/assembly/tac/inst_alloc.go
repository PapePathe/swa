package tac

import "swahili/lang/ast"

// InstAlloc reserves a named stack slot for a variable of type T.
type InstAlloc struct {
	T    ast.Type
	Name string
}

var _ AsmOp = (*InstAlloc)(nil)

func (i *InstAlloc) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstAlloc(i)
}

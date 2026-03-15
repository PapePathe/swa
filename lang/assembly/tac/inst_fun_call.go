package tac

import "fmt"

// InstFunCall emits a call to a named symbol (e.g. "@printf").
// Arguments must have been pushed via preceding InstFunCallArg ops.
type InstFunCall struct {
	Symbol string
}

var _ AsmOp = (*InstFunCall)(nil)

func (i *InstFunCall) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstFunCall(i)
}

func (i *InstFunCall) String() string {
	return fmt.Sprintf("fn call $%s ", i.Symbol)
}

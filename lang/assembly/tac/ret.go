package tac

import "fmt"

type Ret struct {
	Val InstArg
}

var _ AsmOp = (*Ret)(nil)

func (r *Ret) Gen(g AssemblyOpGenerator) error {
	return g.VisitReturn(r)
}

func (i *Ret) String() string {
	return fmt.Sprintf("ret %s ", i.Val.InstructionArg())
}

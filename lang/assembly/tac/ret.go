package tac

type Ret struct {
	Val InstArg
}

func (r *Ret) Gen(g AssemblyOpGenerator) error {
	return g.VisitReturn(r)
}

var _ AsmOp = (*Ret)(nil)

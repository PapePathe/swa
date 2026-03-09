package tac

type InstMul struct {
	Left  InstArg
	Right InstArg
	Width int
}

var _ AsmOp = (*InstMul)(nil)

func (i *InstMul) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstMul(i)
}

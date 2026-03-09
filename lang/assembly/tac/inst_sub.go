package tac

type InstSub struct {
	Left  InstArg
	Right InstArg
	Width int
}

var _ AsmOp = (*InstSub)(nil)

func (i *InstSub) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstSub(i)
}

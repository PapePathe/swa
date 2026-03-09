package tac

type InstAdd struct {
	Left  InstArg
	Right InstArg
	Width int
}

var _ AsmOp = (*InstAdd)(nil)

func (i *InstAdd) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstAdd(i)
}

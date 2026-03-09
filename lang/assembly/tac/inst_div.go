package tac

type InstDiv struct {
	Left  InstArg
	Right InstArg
	Width int
}

var _ AsmOp = (*InstDiv)(nil)

func (i *InstDiv) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstDiv(i)
}

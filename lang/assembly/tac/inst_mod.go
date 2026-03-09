package tac

type InstMod struct {
	Left  InstArg
	Right InstArg
	Width int
}

var _ AsmOp = (*InstMod)(nil)

func (i *InstMod) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstMod(i)
}

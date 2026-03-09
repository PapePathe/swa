package tac

type Label struct {
	Name string
	Ops  []AsmOp
}

var _ InstArg = (*Label)(nil)
var _ AsmOp = (*Label)(nil)

func (l *Label) InstructionArg() string {
	return l.Name
}

func (l *Label) Gen(g AssemblyOpGenerator) error {
	return g.VisitLabel(l)
}

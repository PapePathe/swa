package tac

// InstGlobal declares a read-only string constant in the .rodata section.
// ID is the 0-based index used to generate its label (.Ln).
type InstGlobal struct {
	ID    uint32
	Value string
}

var _ AsmOp = (*InstGlobal)(nil)

func (i *InstGlobal) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstGlobal(i)
}

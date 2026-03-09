package tac

// InstWrite stores Src into the memory location referred to by Dst.
type InstWrite struct {
	Src   InstArg
	Dst   InstArg
	Width int
}

var _ AsmOp = (*InstWrite)(nil)

func (i *InstWrite) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstWrite(i)
}

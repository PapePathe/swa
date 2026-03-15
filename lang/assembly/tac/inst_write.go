package tac

import "fmt"

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

func (i *InstWrite) String() string {
	return fmt.Sprintf("write %s, %s ", i.Src.InstructionArg(), i.Dst.InstructionArg())
}

type InstWriteFloat struct {
	Src   InstArg
	Dst   InstArg
	Width int
}

var _ AsmOp = (*InstWriteFloat)(nil)

func (i *InstWriteFloat) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstWriteFloat(i)
}

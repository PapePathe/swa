package tac

type InstArg interface {
	InstructionArg() string
}

type Inst struct {
	Operation Op
	ArgOne    InstArg
	ArgTwo    InstArg
}

func (i *Inst) Gen(g AssemblyOpGenerator) error {
	return nil
}

var _ AsmOp = (*Inst)(nil)

type AsmOp interface {
	Gen(g AssemblyOpGenerator) error
}

type AssemblyOpGenerator interface {
	VisitBoolVal(node *BoolVal) error
	VisitNumber32Val(node *Number32Val) error
	VisitNumber64Val(node *Number64Val) error
	VisitProc(node *Proc) error
	VisitLabel(node *Label) error
	VisitReturn(node *Ret) error
	VisitTriple(node *Triple) error
}

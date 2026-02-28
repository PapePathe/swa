package generator

import (
	"fmt"
	"strings"
	"swahili/lang/assembly/tac"
)

type Linux64AsmGen struct {
	builder strings.Builder
}

var _ tac.AssemblyOpGenerator = (*Linux64AsmGen)(nil)

func NewLinux64AsmGen() *Linux64AsmGen {
	return &Linux64AsmGen{
		builder: strings.Builder{},
	}
}

func (l *Linux64AsmGen) VisitTriple(node *tac.Triple) error {
	l.builder.WriteString("  .text\n")
	l.builder.WriteString("  .file \"test.swa\"\n")

	for _, proc := range node.Procs {
		err := proc.Gen(l)
		if err != nil {
			return err
		}
	}

	err := node.Main.Gen(l)
	if err != nil {
		return err
	}

	return nil
}

func (l *Linux64AsmGen) Code() string {
	return l.builder.String()
}

func (l *Linux64AsmGen) VisitBoolVal(node *tac.BoolVal) error {
	panic("unimplemented")
}

func (l *Linux64AsmGen) VisitLabel(node *tac.Label) error {
	fmt.Fprintf(&l.builder, "  %s:\n", node.Name)

	for _, op := range node.Ops {
		err := op.Gen(l)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Linux64AsmGen) VisitNumber32Val(node *tac.Number32Val) error {
	panic("unimplemented")
}

func (l *Linux64AsmGen) VisitNumber64Val(node *tac.Number64Val) error {
	panic("unimplemented")
}

func (l *Linux64AsmGen) VisitProc(node *tac.Proc) error {
	fmt.Fprintf(&l.builder, "  .globl %s\n", node.Name)
	fmt.Fprintf(&l.builder, "  .type %s, @function\n", node.Name)
	fmt.Fprintf(&l.builder, "%s:\n", node.Name)

	for _, lab := range node.Labels {
		err := lab.Gen(l)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Linux64AsmGen) VisitReturn(node *tac.Ret) error {
	fmt.Fprintf(&l.builder, "    movl %s, %%eax\n", node.Val.InstructionArg())
	l.builder.WriteString("    retq\n")

	return nil
}

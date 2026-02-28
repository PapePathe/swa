package tac

import (
	"fmt"
	"strings"
)

type Label struct {
	Name  string
	Insts []Inst
	Ops   []AsmOp
}

var _ InstArg = (*Label)(nil)
var _ AsmOp = (*Label)(nil)

func (l *Label) InstructionArg() string {
	return l.Name
}

func (l *Label) Gen(g AssemblyOpGenerator) error {
	return g.VisitLabel(l)
}

func (l *Label) CodeGen(builder *strings.Builder, isdefault bool) {
	for _, inst := range l.Insts {
		fmt.Printf("%+v\n", inst)
	}
}

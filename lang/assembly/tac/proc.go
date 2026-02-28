package tac

import (
	"fmt"
	"swahili/lang/ast"
)

type Proc struct {
	Labels       []*Label
	Name         string
	Ret          ast.Type
	Args         []ast.FuncArg
	Table        map[string]InstID
	currentLabel *Label
	lmap         map[string]int
}

var _ AsmOp = (*Proc)(nil)

func (p *Proc) Gen(g AssemblyOpGenerator) error {
	return g.VisitProc(p)
}

func (p *Proc) Append(i Inst) {
	p.currentLabel.Insts = append(p.currentLabel.Insts, i)
}

func (p *Proc) AppendOp(op AsmOp) {
	p.currentLabel.Ops = append(p.currentLabel.Ops, op)
}

func (p *Proc) addLabel(name string) *Label {
	id, ok := p.lmap[name]
	if !ok {
		p.lmap[name] = 0
		id = 0
	}

	p.lmap[name]++

	newname := fmt.Sprintf("%s-%d", name, id)

	l := &Label{Name: newname, Insts: []Inst{}}

	p.Labels = append(p.Labels, l)

	return l
}

func (p *Proc) setCurrentLabel(l *Label) {
	if l == nil {
		panic("Developer Error label is nil")
	}

	p.currentLabel = l
}

func (p Proc) LastInstID() InstID {
	id := uint32(len(p.currentLabel.Insts) - 1)

	return InstID{id: id}
}

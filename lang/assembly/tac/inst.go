package tac

import (
	"fmt"
	"swahili/lang/ast"
)

type InstArg interface {
	InstructionArg() string
}

type BoolVal struct {
	value string
}

var _ InstArg = (*BoolVal)(nil)

func (b BoolVal) InstructionArg() string {
	return b.value
}

type InstID struct {
	id uint32
}

var _ InstArg = (*InstID)(nil)

func (i InstID) InstructionArg() string {
	return fmt.Sprintf("(%d)", i.id)
}

type GlobalId struct {
	id uint32
}

var _ InstArg = (*GlobalId)(nil)

func (i GlobalId) InstructionArg() string {
	return fmt.Sprintf("@(%d)", i.id)
}

type TypeID struct {
	T ast.Type
}

var _ InstArg = (*TypeID)(nil)

func (t *TypeID) InstructionArg() string {
	return t.T.String()
}

type GlobalString struct {
	Value  string
	Length int
}

var _ InstArg = (*GlobalString)(nil)

func (g *GlobalString) InstructionArg() string {
	return fmt.Sprintf("[%d * byte] = `%q`", g.Length, g.Value)
}

type JumpCond struct {
	Success string
	Failure string
}

var _ InstArg = (*JumpCond)(nil)

func (j *JumpCond) InstructionArg() string {
	return fmt.Sprintf("%s , %s", j.Success, j.Failure)
}

type Inst struct {
	Operation Op
	ArgOne    InstArg
	ArgTwo    InstArg
}

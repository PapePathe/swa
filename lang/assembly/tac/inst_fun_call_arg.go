package tac

import "fmt"

// InstFunCallArg stages one argument for the next InstFunCall.
// Backends accumulate these and assign them to calling-convention
// registers (or the stack) when the call is emitted.
type InstFunCallArg struct {
	Val   InstArg
	Width int
}

var _ AsmOp = (*InstFunCallArg)(nil)

func (i *InstFunCallArg) Gen(g AssemblyOpGenerator) error {
	return g.VisitInstFunCallArg(i)
}

func (i *InstFunCallArg) String() string {
	return fmt.Sprintf("fn arg %s ", i.Val.InstructionArg())
}

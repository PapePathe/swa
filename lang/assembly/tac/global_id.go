package tac

import "fmt"

type GlobalId struct {
	id uint32
}

var _ InstArg = (*GlobalId)(nil)

func (i GlobalId) InstructionArg() string {
	return fmt.Sprintf("@(%d)", i.id)
}

// ID returns the 0-based index of this global, used to generate label names (.L0, .L1, …).
func (i GlobalId) ID() uint32 {
	return i.id
}

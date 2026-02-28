package tac

import "fmt"

type GlobalId struct {
	id uint32
}

var _ InstArg = (*GlobalId)(nil)

func (i GlobalId) InstructionArg() string {
	return fmt.Sprintf("@(%d)", i.id)
}

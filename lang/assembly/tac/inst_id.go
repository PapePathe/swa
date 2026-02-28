package tac

import "fmt"

type InstID struct {
	id uint32
}

var _ InstArg = (*InstID)(nil)

func (i InstID) InstructionArg() string {
	return fmt.Sprintf("(%d)", i.id)
}

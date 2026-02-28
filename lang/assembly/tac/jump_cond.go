package tac

import "fmt"

type JumpCond struct {
	Success string
	Failure string
}

var _ InstArg = (*JumpCond)(nil)

func (j *JumpCond) InstructionArg() string {
	return fmt.Sprintf("%s , %s", j.Success, j.Failure)
}

package tac

type BoolVal struct {
	value string
}

var _ InstArg = (*BoolVal)(nil)

func (b BoolVal) InstructionArg() string {
	return b.value
}

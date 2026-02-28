package tac

type SymbolVal struct {
	value string
}

var _ InstArg = (*SymbolVal)(nil)

func (s *SymbolVal) InstructionArg() string {
	return s.value
}

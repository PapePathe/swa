package tac

import "fmt"

type GlobalString struct {
	Value  string
	Length int
}

var _ InstArg = (*GlobalString)(nil)

func (g *GlobalString) InstructionArg() string {
	return fmt.Sprintf("[%d * byte] = `%q`", g.Length, g.Value)
}

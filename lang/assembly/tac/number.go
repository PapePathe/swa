package tac

import "fmt"

type Number64Val struct {
	value int64
}

func (n *Number64Val) InstructionArg() string {
	return fmt.Sprintf("$%d", n.value)
}

var _ InstArg = (*Number64Val)(nil)

type Number32Val struct {
	value int
}

var _ InstArg = (*Number32Val)(nil)

func (n *Number32Val) InstructionArg() string {
	return fmt.Sprintf("$%d", n.value)
}

type Float64Val struct {
	value float64
}

var _ InstArg = (*Float64Val)(nil)

func (f *Float64Val) InstructionArg() string {
	return fmt.Sprintf("%v", f.value)
}

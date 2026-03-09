package tac

import "fmt"

// OpRef is an InstArg that refers to the result of a previously emitted
// AsmOp by its 0-based index within the current label's Ops slice.
// Backends use this to locate the stack slot that holds a computed value.
type OpRef struct {
	idx int
}

var _ InstArg = (*OpRef)(nil)

// NewOpRef creates an OpRef for the given op index.
func NewOpRef(idx int) *OpRef {
	return &OpRef{idx: idx}
}

// Index returns the op index this ref points to.
func (o *OpRef) Index() int {
	return o.idx
}

func (o *OpRef) InstructionArg() string {
	return fmt.Sprintf("[op%d]", o.idx)
}

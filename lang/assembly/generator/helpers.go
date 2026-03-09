package generator

import (
	"swahili/lang/assembly/tac"
)

// sysVArgRegs is the set of registers used for integer arguments in the
// SysV AMD64 ABI (Linux, macOS, BSDs).
var sysVArgRegs = []string{"%rdi", "%rsi", "%rdx", "%rcx", "%r8", "%r9"}

// preScanProc does a pass over the procedure's instructions to pre-allocate
// stack slots for all named variables, ensuring the stack frame size is
// known before emitting the function prologue.
func preScanProc(alloc *RegAlloc, proc *tac.Proc) {
	for _, lab := range proc.Labels {
		for _, op := range lab.Ops {
			if ia, ok := op.(*tac.InstAlloc); ok {
				alloc.SlotByName(ia.Name)
			}
		}
	}
}

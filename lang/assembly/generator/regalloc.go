package generator

import (
	"swahili/lang/assembly/tac"
)

// RegAlloc is a naïve stack-spill allocator for x86-64 backends.
// Every named variable and every op-result gets its own 8-byte slot
// addressed as negative_offset(%rbp). No physical register is kept
// live across instruction boundaries.
type RegAlloc struct {
	namedSlots map[string]int // variable name → offset
	opSlots    map[int]int    // op index in Label.Ops → offset
	next       int            // next free offset, grows downward from -8
}

func NewRegAlloc() *RegAlloc {
	return &RegAlloc{
		namedSlots: map[string]int{},
		opSlots:    map[int]int{},
		next:       -8,
	}
}

func (r *RegAlloc) alloc() int {
	off := r.next
	r.next -= 8
	return off
}

// SlotByName returns (allocating on first use) the %rbp offset for a named variable.
func (r *RegAlloc) SlotByName(name string) int {
	if off, ok := r.namedSlots[name]; ok {
		return off
	}
	off := r.alloc()
	r.namedSlots[name] = off
	return off
}

// SlotForRef returns (allocating on first use) the %rbp offset for an op result.
func (r *RegAlloc) SlotForRef(ref *tac.OpRef) int {
	idx := ref.Index()
	if off, ok := r.opSlots[idx]; ok {
		return off
	}
	off := r.alloc()
	r.opSlots[idx] = off
	return off
}

// SlotForArg resolves an InstArg to its %rbp offset.
func (r *RegAlloc) SlotForArg(arg tac.InstArg) int {
	switch a := arg.(type) {
	case *tac.OpRef:
		return r.SlotForRef(a)
	case *tac.SymbolVal:
		return r.SlotByName(a.InstructionArg())
	case *tac.GlobalId:
		// Globals (strings) are addressed via rip-relative leaq, not stack slots.
		panic("RegAlloc.SlotForArg: GlobalId should be handled by caller (leaq)")
	default:
		panic("RegAlloc.SlotForArg: cannot allocate slot for " + arg.InstructionArg())
	}
}

// StackSize returns the number of bytes to subtract from %rsp in the prologue.
// Always aligned to 16 bytes for the SysV and Windows x64 ABI.
func (r *RegAlloc) StackSize() int {
	used := -r.next - 8
	if used <= 0 {
		return 0
	}
	// Round up to 16-byte alignment.
	if used%16 != 0 {
		used += 16 - (used % 16)
	}
	return used
}

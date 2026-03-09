package tac

type InstArg interface {
	InstructionArg() string
}

// AsmOp is implemented by every typed instruction and by Proc/Label/Ret.
type AsmOp interface {
	Gen(g AssemblyOpGenerator) error
}

// AssemblyOpGenerator is the visitor interface that every backend must satisfy.
// One Visit method per instruction kind gives backends full type-safe dispatch
// and room for per-instruction local optimisations.
type AssemblyOpGenerator interface {
	// Structural
	VisitTriple(node *Triple) error
	VisitProc(node *Proc) error
	VisitLabel(node *Label) error
	VisitReturn(node *Ret) error

	// Arithmetic
	VisitInstAdd(node *InstAdd) error
	VisitInstSub(noder *InstSub) error
	VisitInstMul(node *InstMul) error
	VisitInstDiv(node *InstDiv) error
	VisitInstMod(node *InstMod) error

	// Memory
	VisitInstAlloc(node *InstAlloc) error
	VisitInstWrite(node *InstWrite) error

	// Calls
	VisitInstFunCall(node *InstFunCall) error
	VisitInstFunCallArg(node *InstFunCallArg) error

	// Data
	VisitInstGlobal(node *InstGlobal) error

	// Leaf values (used when a backend needs to materialise a bare value)
	VisitBoolVal(node *BoolVal) error
	VisitNumber32Val(node *Number32Val) error
	VisitNumber64Val(node *Number64Val) error
}

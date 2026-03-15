package generator

import (
	"fmt"
	"strings"
	"swahili/lang/assembly/tac"
)

// Linux64AsmGen emits x86-64 AT&T-syntax assembly for Linux (ELF, SysV ABI).
//
// Register allocation strategy: naïve stack-spill.
// Every TAC value (named variable, op result, global reference) gets a
// private 8-byte slot addressed as offset(%rbp). No register is live across
// instruction boundaries. The allocator is reset fresh for each Proc/main.
//
// Calling convention (SysV AMD64):
//
//	Integer args 1-6: %rdi %rsi %rdx %rcx %r8 %r9
//	Integer return:   %rax  (we use %eax for 32-bit values)
type Linux64AsmGen struct {
	builder      strings.Builder
	alloc        *RegAlloc
	currentOpIdx int // index of the instruction currently being visited
	// pending args for the next function call, in order
	pendingArgs []struct {
		Arg   tac.InstArg
		Width int
	}
}

var _ tac.AssemblyOpGenerator = (*Linux64AsmGen)(nil)

func NewLinux64AsmGen() *Linux64AsmGen {
	return &Linux64AsmGen{}
}

// -----------------------------------------------------------------------
// Structural visitors
// -----------------------------------------------------------------------

func (l *Linux64AsmGen) VisitTriple(node *tac.Triple) error {
	// Global read-only data (.rodata section for string constants)
	if len(node.GlobalOps) > 0 {
		l.builder.WriteString("  .section .rodata\n")
		for _, op := range node.GlobalOps {
			if err := op.Gen(l); err != nil {
				return err
			}
		}
	}

	l.builder.WriteString("  .text\n")
	l.builder.WriteString("  .file \"swa\"\n")

	for _, proc := range node.Procs {
		if err := proc.Gen(l); err != nil {
			return err
		}
	}

	if node.Main != nil {
		if err := node.Main.Gen(l); err != nil {
			return err
		}
	}

	return nil
}

func (l *Linux64AsmGen) VisitProc(node *tac.Proc) error {
	l.alloc = NewRegAlloc()

	// Pre-scan to calculate stack size before emitting prologue.
	// We do a dry-run allocation pass by visiting allocs.
	preScanProc(l.alloc, node)

	fmt.Fprintf(&l.builder, "  .globl %s\n", node.Name)
	fmt.Fprintf(&l.builder, "  .type %s, @function\n", node.Name)
	fmt.Fprintf(&l.builder, "%s:\n", node.Name)

	l.emitPrologue()

	for _, lab := range node.Labels {
		if err := lab.Gen(l); err != nil {
			return err
		}
	}

	return nil
}

func (l *Linux64AsmGen) VisitLabel(node *tac.Label) error {
	if node.Name != "default" {
		fmt.Fprintf(&l.builder, ".%s:\n", node.Name)
	}

	for i, op := range node.Ops {
		l.currentOpIdx = i
		if err := op.Gen(l); err != nil {
			return err
		}
	}

	return nil
}

func (l *Linux64AsmGen) VisitReturn(node *tac.Ret) error {
	// Return value usually goes in %rax/%eax. For now, we'll use %rax to be safe
	// if it's 64-bit, but the SysV ABI specifies %rax for 64-bit and %eax for 32-bit.
	// Since %eax is the lower half of %rax, using %rax is generally fine if the
	// caller knows what to expect.
	l.loadIntoReg(node.Val, "%rax", 64)
	l.emitEpilogue()

	return nil
}

// -----------------------------------------------------------------------
// Arithmetic visitors — all 32-bit integer for now
// -----------------------------------------------------------------------

func (l *Linux64AsmGen) VisitInstAdd(node *tac.InstAdd) error {
	ref := l.currentOpRef()
	regA, regC, suffix := regAndSuffix(node.Width)

	l.loadIntoReg(node.Left, regA, node.Width)
	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    add%s %s, %s\n", suffix, regC, regA)
	l.storeFromReg(regA, ref, node.Width)

	return nil
}

func (l *Linux64AsmGen) VisitInstSub(node *tac.InstSub) error {
	ref := l.currentOpRef()
	regA, regC, suffix := regAndSuffix(node.Width)

	l.loadIntoReg(node.Left, regA, node.Width)
	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    sub%s %s, %s\n", suffix, regC, regA)
	l.storeFromReg(regA, ref, node.Width)

	return nil
}

func (l *Linux64AsmGen) VisitInstMul(node *tac.InstMul) error {
	ref := l.currentOpRef()
	regA, regC, suffix := regAndSuffix(node.Width)

	l.loadIntoReg(node.Left, regA, node.Width)
	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    imul%s %s, %s\n", suffix, regC, regA)
	l.storeFromReg(regA, ref, node.Width)

	return nil
}
func (l *Linux64AsmGen) VisitInstDiv(node *tac.InstDiv) error {
	ref := l.currentOpRef()

	regA, regC, suffix := regAndSuffix(node.Width)
	if node.Width == 64 {
		l.loadIntoReg(node.Left, regA, 64)
		l.builder.WriteString("    cqto\n") // sign-extend rax → rdx:rax
	} else {
		l.loadIntoReg(node.Left, regA, 32)
		l.builder.WriteString("    cltd\n") // sign-extend eax → edx:eax
	}

	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    idiv%s %s\n", suffix, regC)
	l.storeFromReg(regA, ref, node.Width)

	return nil
}

func (l *Linux64AsmGen) VisitInstMod(node *tac.InstMod) error {
	ref := l.currentOpRef()
	regA, regC, regD := "%eax", "%ecx", "%edx"
	suffix := "l"
	if node.Width == 64 {
		regA, regC, regD = "%rax", "%rcx", "%rdx"
		suffix = "q"
		l.loadIntoReg(node.Left, regA, 64)
		l.builder.WriteString("    cqto\n")
	} else {
		l.loadIntoReg(node.Left, regA, 32)
		l.builder.WriteString("    cltd\n")
	}

	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    idiv%s %s\n", suffix, regC)

	l.storeFromReg(regD, ref, node.Width)

	return nil
}

func regAndSuffix(size int) (string, string, string) {
	regA, regC, suffix := "", "", ""
	switch size {
	case 32:
		regA, regC, suffix = "%eax", "%ecx", "l"
	case 64:
		regA, regC, suffix = "%rax", "%rcx", "q"
	default:
		panic(fmt.Sprintf("unsupported register size %d", size))
	}

	return regA, regC, suffix
}

// -----------------------------------------------------------------------
// Memory visitors
// -----------------------------------------------------------------------

func (l *Linux64AsmGen) VisitInstAlloc(node *tac.InstAlloc) error {
	// Reserve a slot by name. The actual sub %rsp was done in the prologue.
	l.alloc.SlotByName(node.Name)

	return nil
}

func (l *Linux64AsmGen) VisitInstWrite(node *tac.InstWrite) error {
	reg := "%eax"
	suffix := "l"
	if node.Width == 64 {
		reg = "%rax"
		suffix = "q"
	}
	l.loadIntoReg(node.Src, reg, node.Width)
	off := l.alloc.SlotByName(node.Dst.InstructionArg())
	fmt.Fprintf(&l.builder, "    mov%s %s, %d(%%rbp)\n", suffix, reg, off)

	return nil
}

func (l *Linux64AsmGen) VisitInstWriteFloat(node *tac.InstWriteFloat) error {
	reg := "%xmm0"
	suffix := "l"
	if node.Width == 64 {
		reg = "%rax"
		suffix = "q"
	}
	l.loadIntoReg(node.Src, reg, node.Width)
	off := l.alloc.SlotByName(node.Dst.InstructionArg())
	fmt.Fprintf(&l.builder, "    mov%s %s, %d(%%rbp)\n", suffix, reg, off)

	return nil
}

// -----------------------------------------------------------------------
// Call visitors (printf only from main, SysV ABI)
// -----------------------------------------------------------------------

func (l *Linux64AsmGen) VisitInstFunCallArg(node *tac.InstFunCallArg) error {
	// Accumulate args; assign registers when we see the call.
	// We might need width info here too, but for now we store the arg itself.
	l.pendingArgs = append(l.pendingArgs, struct {
		Arg   tac.InstArg
		Width int
	}{Arg: node.Val, Width: node.Width})
	return nil
}

func (l *Linux64AsmGen) VisitInstFunCall(node *tac.InstFunCall) error {
	for i, argEntry := range l.pendingArgs {
		if i >= len(sysVArgRegs) {
			return fmt.Errorf("too many arguments")
		}
		reg := sysVArgRegs[i]
		arg := argEntry.Arg
		width := argEntry.Width

		switch a := arg.(type) {
		case *tac.GlobalId:
			// Load address of the global string into the 64-bit arg register.
			fmt.Fprintf(&l.builder, "    leaq .L%d(%%rip), %s\n", a.ID(), reg)
		default:
			// Load from stack slot.
			off := l.alloc.SlotForArg(arg)
			if width == 64 {
				fmt.Fprintf(&l.builder, "    movq %d(%%rbp), %s\n", off, reg)
			} else {
				// Sign-extend 32-bit to 64-bit for the argument register.
				fmt.Fprintf(&l.builder, "    movslq %d(%%rbp), %s\n", off, reg)
			}
		}
	}

	l.pendingArgs = nil

	// For variadic functions (printf) %al must hold the number of float args.
	l.builder.WriteString("    movl $0, %eax\n")
	fmt.Fprintf(&l.builder, "    call %s\n", node.Symbol)

	return nil
}

// -----------------------------------------------------------------------
// Data visitors
// -----------------------------------------------------------------------

func (l *Linux64AsmGen) VisitInstGlobal(node *tac.InstGlobal) error {
	fmt.Fprintf(&l.builder, ".L%d:\n", node.ID)
	fmt.Fprintf(&l.builder, "  .string %q\n", node.Value)

	return nil
}

// -----------------------------------------------------------------------
// Leaf value visitors (rarely called directly from outside)
// -----------------------------------------------------------------------

func (l *Linux64AsmGen) VisitBoolVal(node *tac.BoolVal) error         { return nil }
func (l *Linux64AsmGen) VisitNumber32Val(node *tac.Number32Val) error { return nil }
func (l *Linux64AsmGen) VisitNumber64Val(node *tac.Number64Val) error { return nil }

// -----------------------------------------------------------------------
// Code() returns the accumulated assembly text.
// -----------------------------------------------------------------------

func (l *Linux64AsmGen) Code() string {
	return l.builder.String()
}

// -----------------------------------------------------------------------
// Internal helpers
// -----------------------------------------------------------------------

func (l *Linux64AsmGen) emitPrologue() {
	l.builder.WriteString("    pushq %rbp\n")
	l.builder.WriteString("    movq %rsp, %rbp\n")
	sz := l.alloc.StackSize()
	if sz > 0 {
		fmt.Fprintf(&l.builder, "    subq $%d, %%rsp\n", sz)
	}
}

func (l *Linux64AsmGen) emitEpilogue() {
	l.builder.WriteString("    movq %rbp, %rsp\n")
	l.builder.WriteString("    popq %rbp\n")
	l.builder.WriteString("    retq\n")
}

// loadIntoReg loads arg into destReg (%rax, %eax, …).
// Immediates are moved directly; stack values are loaded via movl/movq.
func (l *Linux64AsmGen) loadIntoReg(arg tac.InstArg, destReg string, width int) {
	suffix := "l"
	if width == 64 {
		suffix = "q"
	}

	switch a := arg.(type) {
	case *tac.Number32Val, *tac.Number64Val:
		fmt.Fprintf(&l.builder, "    mov%s %s, %s\n", suffix, arg.InstructionArg(), destReg)
	case *tac.OpRef:
		off := l.alloc.SlotForRef(a)
		fmt.Fprintf(&l.builder, "    mov%s %d(%%rbp), %s\n", suffix, off, destReg)
	case *tac.SymbolVal:
		off := l.alloc.SlotByName(a.InstructionArg())
		fmt.Fprintf(&l.builder, "    mov%s %d(%%rbp), %s\n", suffix, off, destReg)
	default:
		// Fallback: treat InstructionArg() as an AT&T source operand.
		fmt.Fprintf(&l.builder, "    mov%s %s, %s\n", suffix, arg.InstructionArg(), destReg)
	}
}

// storeFromReg stores the value in srcReg into the stack slot for ref.
func (l *Linux64AsmGen) storeFromReg(srcReg string, ref *tac.OpRef, width int) {
	suffix := "l"
	if width == 64 {
		suffix = "q"
	}
	off := l.alloc.SlotForRef(ref)
	fmt.Fprintf(&l.builder, "    mov%s %s, %d(%%rbp)\n", suffix, srcReg, off)
}

// currentOpRef returns an OpRef pointing to the slot that will hold the
// result of the instruction currently being emitted.
// Must be called at the start of an arithmetic visitor, before appending
// anything, so the index equals the op's own position in the label.
func (l *Linux64AsmGen) currentOpRef() *tac.OpRef {
	return tac.NewOpRef(l.currentOpIdx)
}

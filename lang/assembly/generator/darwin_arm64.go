package generator

import (
	"fmt"
	"strings"
	"swahili/lang/assembly/tac"
)

// NewPlatformAsmGen returns the appropriate assembly generator for the given
// operating system and architecture.
func NewPlatformAsmGen(os, arch string) (tac.AssemblyOpGenerator, error) {
	switch os {
	case "linux", "Gnu/Linux":
		if arch == "amd64" || arch == "X86-64" {
			return NewLinux64AsmGen(), nil
		}
		//	case "darwin":
		//		if arch == "amd64" {
		//			return NewDarwinAmd64AsmGen(), nil
		//		}
		//		if arch == "arm64" {
		//			return NewDarwinArm64AsmGen(), nil
		//		}
		//	case "windows":
		//		if arch == "amd64" {
		//			return NewWindowsAmd64AsmGen(), nil
		//		}
	}

	return nil, fmt.Errorf("unsupported target: %s/%s", os, arch)
}

// DarwinArm64AsmGen emits ARM64 assembly for macOS (Mach-O, AAPCS64 ABI).
type DarwinArm64AsmGen struct {
	builder      strings.Builder
	alloc        *RegAlloc
	currentOpIdx int
	// pending args for the next function call, in order
	pendingArgs []struct {
		Arg   tac.InstArg
		Width int
	}
}

// VisitInstWriteFloat implements [tac.AssemblyOpGenerator].
func (l *DarwinArm64AsmGen) VisitInstWriteFloat(node *tac.InstWriteFloat) error {
	panic("unimplemented")
}

var _ tac.AssemblyOpGenerator = (*DarwinArm64AsmGen)(nil)

func newdarwinarm64asmgen() *DarwinArm64AsmGen {
	return &DarwinArm64AsmGen{}
}

func (l *DarwinArm64AsmGen) VisitTriple(node *tac.Triple) error {
	if len(node.GlobalOps) > 0 {
		l.builder.WriteString("  .section __TEXT,__cstring\n")
		for _, op := range node.GlobalOps {
			if err := op.Gen(l); err != nil {
				return err
			}
		}
	}

	l.builder.WriteString("  .section __TEXT,__text\n")
	l.builder.WriteString("  .align 2\n")

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

func (l *DarwinArm64AsmGen) VisitProc(node *tac.Proc) error {
	l.alloc = NewRegAlloc()
	preScanProc(l.alloc, node)

	name := node.Name
	if name == "main" {
		name = "_main"
	} else {
		name = "_" + name
	}

	fmt.Fprintf(&l.builder, "  .globl %s\n", name)
	fmt.Fprintf(&l.builder, "  .p2align 2\n")
	fmt.Fprintf(&l.builder, "%s:\n", name)

	l.emitPrologue()

	for _, lab := range node.Labels {
		if err := lab.Gen(l); err != nil {
			return err
		}
	}

	return nil
}

func (l *DarwinArm64AsmGen) VisitLabel(node *tac.Label) error {
	if node.Name != "default" {
		fmt.Fprintf(&l.builder, "L%s:\n", node.Name)
	}

	for i, op := range node.Ops {
		l.currentOpIdx = i
		if err := op.Gen(l); err != nil {
			return err
		}
	}

	return nil
}

func (l *DarwinArm64AsmGen) VisitReturn(node *tac.Ret) error {
	l.loadIntoReg(node.Val, "x0", 64)
	l.emitEpilogue()
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstAdd(node *tac.InstAdd) error {
	ref := l.currentOpRef()
	reg9, reg10 := "w9", "w10"
	if node.Width == 64 {
		reg9, reg10 = "x9", "x10"
	}
	l.loadIntoReg(node.Left, reg9, node.Width)
	l.loadIntoReg(node.Right, reg10, node.Width)
	fmt.Fprintf(&l.builder, "    add %s, %s, %s\n", reg9, reg9, reg10)
	l.storeFromReg(reg9, ref, node.Width)
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstSub(node *tac.InstSub) error {
	ref := l.currentOpRef()
	reg9, reg10 := "w9", "w10"
	if node.Width == 64 {
		reg9, reg10 = "x9", "x10"
	}
	l.loadIntoReg(node.Left, reg9, node.Width)
	l.loadIntoReg(node.Right, reg10, node.Width)
	fmt.Fprintf(&l.builder, "    sub %s, %s, %s\n", reg9, reg9, reg10)
	l.storeFromReg(reg9, ref, node.Width)
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstMul(node *tac.InstMul) error {
	ref := l.currentOpRef()
	reg9, reg10 := "w9", "w10"
	if node.Width == 64 {
		reg9, reg10 = "x9", "x10"
	}
	l.loadIntoReg(node.Left, reg9, node.Width)
	l.loadIntoReg(node.Right, reg10, node.Width)
	fmt.Fprintf(&l.builder, "    mul %s, %s, %s\n", reg9, reg9, reg10)
	l.storeFromReg(reg9, ref, node.Width)
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstDiv(node *tac.InstDiv) error {
	ref := l.currentOpRef()
	reg9, reg10 := "w9", "w10"
	if node.Width == 64 {
		reg9, reg10 = "x9", "x10"
	}
	l.loadIntoReg(node.Left, reg9, node.Width)
	l.loadIntoReg(node.Right, reg10, node.Width)
	fmt.Fprintf(&l.builder, "    sdiv %s, %s, %s\n", reg9, reg9, reg10)
	l.storeFromReg(reg9, ref, node.Width)
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstMod(node *tac.InstMod) error {
	ref := l.currentOpRef()
	reg9, reg10, reg11 := "w9", "w10", "w11"
	if node.Width == 64 {
		reg9, reg10, reg11 = "x9", "x10", "x11"
	}
	l.loadIntoReg(node.Left, reg9, node.Width)
	l.loadIntoReg(node.Right, reg10, node.Width)
	fmt.Fprintf(&l.builder, "    sdiv %s, %s, %s\n", reg11, reg9, reg10)
	fmt.Fprintf(&l.builder, "    msub %s, %s, %s, %s\n", reg9, reg11, reg10, reg9)
	l.storeFromReg(reg9, ref, node.Width)
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstAlloc(node *tac.InstAlloc) error {
	l.alloc.SlotByName(node.Name)
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstWrite(node *tac.InstWrite) error {
	reg := "w9"
	if node.Width == 64 {
		reg = "x9"
	}
	l.loadIntoReg(node.Src, reg, node.Width)
	off := l.alloc.SlotByName(node.Dst.InstructionArg())
	fmt.Fprintf(&l.builder, "    str %s, [sp, #%d]\n", reg, off)
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstFunCallArg(node *tac.InstFunCallArg) error {
	l.pendingArgs = append(l.pendingArgs, struct {
		Arg   tac.InstArg
		Width int
	}{Arg: node.Val, Width: node.Width})
	return nil
}

var arm64ArgRegs = []string{"x0", "x1", "x2", "x3", "x4", "x5", "x6", "x7"}

func (l *DarwinArm64AsmGen) VisitInstFunCall(node *tac.InstFunCall) error {
	for i, argEntry := range l.pendingArgs {
		if i >= 8 { // Darwin/ARM64 has 8 arg registers x0-x7.
			return fmt.Errorf("too many arguments")
		}
		destReg := fmt.Sprintf("x%d", i)
		if argEntry.Width == 32 {
			destReg = fmt.Sprintf("w%d", i)
		}

		arg := argEntry.Arg
		switch a := arg.(type) {
		case *tac.GlobalId:
			// ARM64 usually uses adrp/add or adr for globals.
			// Darwin/ARM64 Mach-O:
			fmt.Fprintf(&l.builder, "    adrp x%d, L%d@PAGE\n", i, a.ID())
			fmt.Fprintf(&l.builder, "    add x%d, x%d, L%d@PAGEOFF\n", i, i, a.ID())
		default:
			off := l.alloc.SlotForArg(arg)
			fmt.Fprintf(&l.builder, "    ldr %s, [sp, #%d]\n", destReg, off)
		}
	}

	l.pendingArgs = nil
	fmt.Fprintf(&l.builder, "    bl _%s\n", node.Symbol)
	return nil
}

func (l *DarwinArm64AsmGen) VisitInstGlobal(node *tac.InstGlobal) error {
	fmt.Fprintf(&l.builder, "L%d:\n", node.ID)
	fmt.Fprintf(&l.builder, "  .asciz %q\n", node.Value)
	return nil
}

func (l *DarwinArm64AsmGen) VisitBoolVal(node *tac.BoolVal) error         { return nil }
func (l *DarwinArm64AsmGen) VisitNumber32Val(node *tac.Number32Val) error { return nil }
func (l *DarwinArm64AsmGen) VisitNumber64Val(node *tac.Number64Val) error { return nil }

func (l *DarwinArm64AsmGen) Code() string {
	return l.builder.String()
}

func (l *DarwinArm64AsmGen) emitPrologue() {
	sz := l.alloc.StackSize()
	// ARM64 requires 16-byte aligned stack. NewRegAlloc already does this.
	// We also need to save FP and LR.
	sz += 16
	fmt.Fprintf(&l.builder, "    stp x29, x30, [sp, #-%d]!\n", sz)
	l.builder.WriteString("    mov x29, sp\n")
}

func (l *DarwinArm64AsmGen) emitEpilogue() {
	sz := l.alloc.StackSize() + 16
	l.builder.WriteString("    mov sp, x29\n")
	fmt.Fprintf(&l.builder, "    ldp x29, x30, [sp], #%d\n", sz)
	l.builder.WriteString("    ret\n")
}

func (l *DarwinArm64AsmGen) loadIntoReg(arg tac.InstArg, destReg string, width int) {
	switch a := arg.(type) {
	case *tac.Number32Val, *tac.Number64Val:
		val := strings.TrimPrefix(arg.InstructionArg(), "$")
		fmt.Fprintf(&l.builder, "    mov %s, #%s\n", destReg, val)
	case *tac.OpRef:
		off := l.alloc.SlotForRef(a)
		fmt.Fprintf(&l.builder, "    ldr %s, [sp, #%d]\n", destReg, off)
	case *tac.SymbolVal:
		off := l.alloc.SlotByName(a.InstructionArg())
		fmt.Fprintf(&l.builder, "    ldr %s, [sp, #%d]\n", destReg, off)
	default:
		fmt.Fprintf(&l.builder, "    // Unsupported load for %T\n", arg)
	}
}

func (l *DarwinArm64AsmGen) storeFromReg(srcReg string, ref *tac.OpRef, width int) {
	off := l.alloc.SlotForRef(ref)
	fmt.Fprintf(&l.builder, "    str %s, [sp, #%d]\n", srcReg, off)
}

func (l *DarwinArm64AsmGen) currentOpRef() *tac.OpRef {
	return tac.NewOpRef(l.currentOpIdx)
}

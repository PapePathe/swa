package generator

import (
	"fmt"
	"strings"
	"swahili/lang/assembly/tac"
)

// WindowsAmd64AsmGen emits x86-64 assembly for Windows (PE/COFF, Windows x64 ABI).
type WindowsAmd64AsmGen struct {
	builder      strings.Builder
	alloc        *RegAlloc
	currentOpIdx int
	// pending args for the next function call, in order
	pendingArgs []struct {
		Arg   tac.InstArg
		Width int
	}
}

func (l *WindowsAmd64AsmGen) VisitInstWriteFloat(node *tac.InstWriteFloat) error {
	panic("unimplemented")
}

var _ tac.AssemblyOpGenerator = (*WindowsAmd64AsmGen)(nil)

func NewWindowsAmd64AsmGen() *WindowsAmd64AsmGen {
	return &WindowsAmd64AsmGen{}
}

func (l *WindowsAmd64AsmGen) VisitTriple(node *tac.Triple) error {
	if len(node.GlobalOps) > 0 {
		l.builder.WriteString("  .section .rdata,\"dr\"\n")
		for _, op := range node.GlobalOps {
			if err := op.Gen(l); err != nil {
				return err
			}
		}
	}

	l.builder.WriteString("  .text\n")

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

func (l *WindowsAmd64AsmGen) VisitProc(node *tac.Proc) error {
	l.alloc = NewRegAlloc()
	preScanProc(l.alloc, node)

	name := node.Name

	fmt.Fprintf(&l.builder, "  .globl %s\n", name)
	fmt.Fprintf(&l.builder, "  .def %s; .scl 2; .type 32; .endef\n", name)
	fmt.Fprintf(&l.builder, "%s:\n", name)

	l.emitPrologue()

	for _, lab := range node.Labels {
		if err := lab.Gen(l); err != nil {
			return err
		}
	}

	return nil
}

func (l *WindowsAmd64AsmGen) VisitLabel(node *tac.Label) error {
	if node.Name != "default" {
		fmt.Fprintf(&l.builder, ".L%s:\n", node.Name)
	}

	for i, op := range node.Ops {
		l.currentOpIdx = i
		if err := op.Gen(l); err != nil {
			return err
		}
	}

	return nil
}

func (l *WindowsAmd64AsmGen) VisitReturn(node *tac.Ret) error {
	l.loadIntoReg(node.Val, "%rax", 64)
	l.emitEpilogue()
	return nil
}

func (l *WindowsAmd64AsmGen) VisitInstAdd(node *tac.InstAdd) error {
	ref := l.currentOpRef()
	regA, regC := "%eax", "%ecx"
	suffix := "l"
	if node.Width == 64 {
		regA, regC = "%rax", "%rcx"
		suffix = "q"
	}

	l.loadIntoReg(node.Left, regA, node.Width)
	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    add%s %s, %s\n", suffix, regC, regA)
	l.storeFromReg(regA, ref, node.Width)
	return nil
}

func (l *WindowsAmd64AsmGen) VisitInstSub(node *tac.InstSub) error {
	ref := l.currentOpRef()
	regA, regC := "%eax", "%ecx"
	suffix := "l"
	if node.Width == 64 {
		regA, regC = "%rax", "%rcx"
		suffix = "q"
	}

	l.loadIntoReg(node.Left, regA, node.Width)
	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    sub%s %s, %s\n", suffix, regC, regA)
	l.storeFromReg(regA, ref, node.Width)
	return nil
}

func (l *WindowsAmd64AsmGen) VisitInstMul(node *tac.InstMul) error {
	ref := l.currentOpRef()
	regA, regC := "%eax", "%ecx"
	suffix := "l"
	if node.Width == 64 {
		regA, regC = "%rax", "%rcx"
		suffix = "q"
	}

	l.loadIntoReg(node.Left, regA, node.Width)
	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    imul%s %s, %s\n", suffix, regC, regA)
	l.storeFromReg(regA, ref, node.Width)
	return nil
}

func (l *WindowsAmd64AsmGen) VisitInstDiv(node *tac.InstDiv) error {
	ref := l.currentOpRef()
	regA, regC := "%eax", "%ecx"
	suffix := "l"
	if node.Width == 64 {
		regA, regC = "%rax", "%rcx"
		suffix = "q"
		l.loadIntoReg(node.Left, regA, 64)
		l.builder.WriteString("    cqto\n")
	} else {
		l.loadIntoReg(node.Left, regA, 32)
		l.builder.WriteString("    cltd\n")
	}

	l.loadIntoReg(node.Right, regC, node.Width)
	fmt.Fprintf(&l.builder, "    idiv%s %s\n", suffix, regC)
	l.storeFromReg(regA, ref, node.Width)
	return nil
}

func (l *WindowsAmd64AsmGen) VisitInstMod(node *tac.InstMod) error {
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

func (l *WindowsAmd64AsmGen) VisitInstAlloc(node *tac.InstAlloc) error {
	l.alloc.SlotByName(node.Name)
	return nil
}

func (l *WindowsAmd64AsmGen) VisitInstWrite(node *tac.InstWrite) error {
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

func (l *WindowsAmd64AsmGen) VisitInstFunCallArg(node *tac.InstFunCallArg) error {
	l.pendingArgs = append(l.pendingArgs, struct {
		Arg   tac.InstArg
		Width int
	}{Arg: node.Val, Width: node.Width})
	return nil
}

var winArgRegs = []string{"%rcx", "%rdx", "%r8", "%r9"}

func (l *WindowsAmd64AsmGen) VisitInstFunCall(node *tac.InstFunCall) error {
	for i, argEntry := range l.pendingArgs {
		if i >= len(winArgRegs) {
			return fmt.Errorf("too many arguments")
		}
		reg := winArgRegs[i]
		arg := argEntry.Arg
		width := argEntry.Width

		switch a := arg.(type) {
		case *tac.GlobalId:
			fmt.Fprintf(&l.builder, "    leaq .L%d(%%rip), %s\n", a.ID(), reg)
		default:
			off := l.alloc.SlotForArg(arg)
			if width == 64 {
				fmt.Fprintf(&l.builder, "    movq %d(%%rbp), %s\n", off, reg)
			} else {
				fmt.Fprintf(&l.builder, "    movslq %d(%%rbp), %s\n", off, reg)
			}
		}
	}

	l.pendingArgs = nil
	// Windows x64 ABI does not use %al for variadic functions (it uses a shadow space).
	fmt.Fprintf(&l.builder, "    call %s\n", node.Symbol)
	return nil
}

func (l *WindowsAmd64AsmGen) VisitInstGlobal(node *tac.InstGlobal) error {
	fmt.Fprintf(&l.builder, ".L%d:\n", node.ID)
	fmt.Fprintf(&l.builder, "  .asciz %q\n", node.Value)
	return nil
}

func (l *WindowsAmd64AsmGen) VisitBoolVal(node *tac.BoolVal) error         { return nil }
func (l *WindowsAmd64AsmGen) VisitNumber32Val(node *tac.Number32Val) error { return nil }
func (l *WindowsAmd64AsmGen) VisitNumber64Val(node *tac.Number64Val) error { return nil }

func (l *WindowsAmd64AsmGen) Code() string {
	return l.builder.String()
}

func (l *WindowsAmd64AsmGen) emitPrologue() {
	l.builder.WriteString("    pushq %rbp\n")
	l.builder.WriteString("    movq %rsp, %rbp\n")
	sz := l.alloc.StackSize()
	// Windows x64 requires 32 bytes of shadow space for the callee.
	sz += 32
	if sz%16 != 0 {
		sz += 16 - (sz % 16)
	}
	fmt.Fprintf(&l.builder, "    subq $%d, %%rsp\n", sz)
}

func (l *WindowsAmd64AsmGen) emitEpilogue() {
	l.builder.WriteString("    movq %rbp, %rsp\n") // This takes care of shadow space and local slots.
	l.builder.WriteString("    popq %rbp\n")
	l.builder.WriteString("    ret\n")
}

func (l *WindowsAmd64AsmGen) loadIntoReg(arg tac.InstArg, destReg string, width int) {
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
		fmt.Fprintf(&l.builder, "    mov%s %s, %s\n", suffix, arg.InstructionArg(), destReg)
	}
}

func (l *WindowsAmd64AsmGen) storeFromReg(srcReg string, ref *tac.OpRef, width int) {
	suffix := "l"
	if width == 64 {
		suffix = "q"
	}
	off := l.alloc.SlotForRef(ref)
	fmt.Fprintf(&l.builder, "    mov%s %s, %d(%%rbp)\n", suffix, srcReg, off)
}

func (l *WindowsAmd64AsmGen) currentOpRef() *tac.OpRef {
	return tac.NewOpRef(l.currentOpIdx)
}

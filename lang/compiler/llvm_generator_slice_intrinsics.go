package compiler

import (
	"tinygo.org/x/go-llvm"
)

type ExternalFunc struct {
	lltype llvm.Type
	llval  llvm.Value
}

func (g *LLVMGenerator) getReallocFunc() ExternalFunc {
	name := "realloc"
	f := g.Ctx.Module.NamedFunction(name)
	i64 := g.Ctx.Context.Int64Type()
	ptrI8 := llvm.PointerType(g.Ctx.Context.Int8Type(), 0)

	ftype := llvm.FunctionType(ptrI8, []llvm.Type{ptrI8, i64}, false)
	if f.IsNil() {
		f = llvm.AddFunction(*g.Ctx.Module, name, ftype)
	}

	return ExternalFunc{lltype: ftype, llval: f}
}

package vm

//import (
//	"github.com/llir/llvm/ir"
//	"github.com/llir/llvm/ir/types"
//	"github.com/llir/llvm/ir/value"
//)
//
//type symbolValue struct {
//	v value.Value
//}
//
//type symbolTable struct {
//	parent  *symbolTable
//	symbols map[string]symbolValue
//}
//
//func (st *symbolTable) Lookup(name string) value.Value {
//	return nil
//}
//
//type Instruction interface {
//	CodeGen(m *ir.Module, block *ir.Block, table *symbolTable) error
//}
//type Program struct {
//	st  symbolTable
//	ins []Instruction
//}
//
//type (
//	I32Def struct {
//		v types.IntType
//	}
//	CharArrayDef struct {
//		v types.PointerType
//	}
//	ArrayDef struct {
//		v types.ArrayType
//	}
//	StructDef struct{}
//)
//
//func (receiver I32Def) CodeGen(m *ir.Module, block *ir.Block, table *symbolTable) error {
//	return nil
//}

package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	name, ok := node.Name.(ast.SymbolExpression)
	if !ok {
		return fmt.Errorf("FunctionCallExpression: name is not a symbol")
	}

	err, funcType := g.Ctx.FindFuncSymbol(name.Value)
	if err != nil {
		return err
	}

	funcVal := g.Ctx.Module.NamedFunction(name.Value)
	if funcVal.IsNil() {
		return fmt.Errorf("function %s does not exist", name.Value)
	}

	if funcVal.ParamsCount() != len(node.Args) {
		return fmt.Errorf("function %s expect %d arguments but was given %d", name.Value, funcVal.ParamsCount(), len(node.Args))
	}

	args := []llvm.Value{}

	for i, arg := range node.Args {
		err := arg.Accept(g)
		if err != nil {
			return err
		}

		val := g.getLastResult()
		if val == nil || val.Value == nil {
			return fmt.Errorf("failed to evaluate argument %d", i+1)
		}

		param := funcVal.Params()[i]
		paramType := param.Type()

		// Type checking
		argType := val.Value.Type()
		if val.SymbolTableEntry != nil && val.SymbolTableEntry.Ref != nil {
			argType = val.SymbolTableEntry.Ref.LLVMType
		}

		if argType != paramType {
			err := fmt.Errorf(
				"expected argument of type %s expected but got %s",
				g.formatLLVMType(paramType),
				g.formatLLVMType(argType),
			)

			if paramType.TypeKind() == llvm.IntegerTypeKind && argType.TypeKind() == llvm.PointerTypeKind {
				return err
			}

			if (argType.TypeKind() == llvm.StructTypeKind && paramType.TypeKind() == llvm.PointerTypeKind) ||
				(argType.TypeKind() == llvm.ArrayTypeKind && paramType.TypeKind() == llvm.PointerTypeKind) {
			} else {
				return err
			}
		}

		switch arg.(type) {
		case ast.SymbolExpression:
			if val.SymbolTableEntry.Ref != nil {
				alloca := g.Ctx.Builder.CreateAlloca(val.SymbolTableEntry.Ref.LLVMType, "")
				g.Ctx.Builder.CreateStore(*val.Value, alloca)
				args = append(args, alloca)

				break
			}

			// Pass arrays by reference
			if _, ok := val.SymbolTableEntry.DeclaredType.(ast.ArrayType); ok {
				args = append(args, *val.SymbolTableEntry.Address)
				break
			}

			// For simple types (int, float, etc.), use the loaded value
			// VisitSymbolExpression already loaded the value for us
			args = append(args, *val.Value)

		default:
			args = append(args, *val.Value)
		}
	}

	val := g.Ctx.Builder.CreateCall(*funcType, funcVal, args, "")

	g.setLastResult(&ast.CompilerResult{Value: &val})

	return nil
}

func (g *LLVMGenerator) formatLLVMType(t llvm.Type) string {
	switch t.TypeKind() {
	case llvm.IntegerTypeKind:
		return fmt.Sprintf("IntegerType(%d bits)", t.IntTypeWidth())
	case llvm.FloatTypeKind:
		return "FloatType"
	case llvm.DoubleTypeKind:
		return "DoubleType"
	case llvm.PointerTypeKind:
		return fmt.Sprintf("PointerType(%s)", g.formatLLVMType(t.ElementType()))
	case llvm.ArrayTypeKind:
		return fmt.Sprintf("ArrayType(%s[%d])", g.formatLLVMType(t.ElementType()), t.ArrayLength())
	case llvm.StructTypeKind:
		return "StructType"
	case llvm.VoidTypeKind:
		return "VoidType"
	default:
		if t.TypeKind() == llvm.TypeKind(1) {
			return "Reference"
		}
		return t.String()
	}
}

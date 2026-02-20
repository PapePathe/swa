package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	old := g.logger.Step("FunCallExpr")

	defer g.logger.Restore(old)

	g.Debugf("%s", node.Name)

	name, ok := node.Name.(*ast.SymbolExpression)
	if !ok {
		key := "LLVMGenerator.VisitFunctionCall.NameIsNotASymbol"

		return g.Ctx.Dialect.Error(key)
	}

	err, funcType := g.Ctx.FindFuncSymbol(name.Value)
	if err != nil {
		return err
	}

	funcVal := g.Ctx.Module.NamedFunction(name.Value)
	if funcVal.IsNil() {
		key := "LLVMGenerator.VisitFunctionCall.DoesNotExist"

		return g.Ctx.Dialect.Error(key, name.Value)
	}

	if funcVal.ParamsCount() != len(node.Args) {
		key := "LLVMGenerator.VisitFunctionCall.ArgsAndParamsCountAreDifferent"

		return g.Ctx.Dialect.Error(key, name.Value, funcVal.ParamsCount(), len(node.Args))
	}

	node.SwaType = funcType.meta.ReturnType
	args := []llvm.Value{}

	for i, arg := range node.Args {
		err := arg.Accept(g)
		if err != nil {
			return err
		}

		val := g.getLastResult()

		g.Debugf("%+v", val)

		if val == nil || val.Value == nil {
			key := "LLVMGenerator.VisitFunctionCall.FailedToEvaluate"

			return g.Ctx.Dialect.Error(key, i+1)
		}

		param := funcVal.Params()[i]
		paramType := param.Type()

		g.Debugf("Symbol Table entry: %+v", val.SymbolTableEntry)
		g.Debugf("array symbol Table entry: %+v", val.ArraySymbolTableEntry)

		// TODO this should be moved to the type checking pass
		argType := val.Value.Type()
		if val.SymbolTableEntry != nil && val.SymbolTableEntry.Ref != nil {
			if val.StuctPropertyValueType != nil {
				argType = *val.StuctPropertyValueType
			} else {
				argType = val.SymbolTableEntry.Ref.LLVMType
			}
		}

		if val.ArraySymbolTableEntry != nil {
			argType = val.ArraySymbolTableEntry.UnderlyingType

			if argType.TypeKind() == llvm.StructTypeKind {
				if val.StuctPropertyValueType == nil {
					key := "LLVMGenerator.VisitFunctionCall.StructPropertyValueTypeIsNil"

					return g.Ctx.Dialect.Error(key)
				}

				argType = *val.StuctPropertyValueType
			}
		}

		if argType != paramType {
			key := "LLVMGenerator.VisitFunctionCall.UnexpectedArgumentType"
			err := g.Ctx.Dialect.Error(key, g.formatLLVMType(paramType), g.formatLLVMType(argType))

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
		case *ast.MemberExpression, *ast.ArrayAccessExpression,
			*ast.ArrayOfStructsAccessExpression:
			load := g.Ctx.Builder.CreateLoad(argType, *val.Value, "")
			args = append(args, load)

		case *ast.SymbolExpression:
			if val.SymbolTableEntry.Ref != nil {
				alloca := g.Ctx.Builder.CreateAlloca(val.SymbolTableEntry.Ref.LLVMType, "")
				g.Ctx.Builder.CreateStore(*val.Value, alloca)
				args = append(args, alloca)

				break
			}

			if _, ok := val.SymbolTableEntry.DeclaredType.(ast.ArrayType); ok {
				args = append(args, *val.SymbolTableEntry.Address)

				break
			}

			args = append(args, *val.Value)

		case *ast.FloatExpression, *ast.NumberExpression,
			*ast.StringExpression, *ast.BinaryExpression,
			*ast.FunctionCallExpression, *ast.PrefixExpression:
			args = append(args, *val.Value)
		default:
			key := "LLVMGenerator.VisitFunctionCall.UnsupportedType"

			return g.Ctx.Dialect.Error(key, arg)
		}
	}

	val := g.Ctx.Builder.CreateCall(funcType.lltype, funcVal, args, "")

	g.setLastResult(&CompilerResult{
		Value:   &val,
		SwaType: node.SwaType,
	})

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
		return "Cannot stringify type "
	}
}

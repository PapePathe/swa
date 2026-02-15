package compiler

import (
	"fmt"
	"reflect"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

type PrintableValueExtractor func(g *LLVMGenerator, res *CompilerResult) llvm.Value

// VisitPrintStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitPrintStatement(node *ast.PrintStatetement) error {
	old := g.logger.Step("PrintStmt")

	defer g.logger.Restore(old)

	g.Debugf("%s", node.Values)

	printableValues := []llvm.Value{}
	results := []*CompilerResult{}

	for _, expr := range node.Values {
		err := expr.Accept(g)
		if err != nil {
			return err
		}

		res := g.getLastResult()
		results = append(results, res)

		extractor, ok := printableValueExtractors[reflect.TypeOf(expr)]
		if !ok {
			// TODO this is a developer error
			// We should inform the user that it's not his fault
			g.NotImplemented(fmt.Sprintf("VisitPrintStatement unimplemented for %T", expr))

			continue
		}

		printableValues = append(printableValues, extractor(g, res))
	}

	if len(printableValues) == 0 {
		return nil
	}

	isLegacy := false
	if firstArg, ok := node.Values[0].(*ast.StringExpression); ok {
		specifierCount := 0
		for i := 0; i < len(firstArg.Value); i++ {
			if firstArg.Value[i] == '%' {
				if i+1 < len(firstArg.Value) && firstArg.Value[i+1] == '%' {
					i++
					continue
				}
				specifierCount++
			}
		}
		if specifierCount > 0 && len(node.Values) > 1 {
			isLegacy = true
		}
	}

	var finalArgs []llvm.Value
	if isLegacy {
		finalArgs = printableValues
	} else {
		formatStr := ""
		for _, res := range results {
			formatStr += g.getFormatSpecifier(res)
		}

		fmtPtr := g.Ctx.Builder.CreateGlobalStringPtr(formatStr, "print.format")
		finalArgs = append([]llvm.Value{fmtPtr}, printableValues...)
	}

	g.Ctx.Builder.CreateCall(
		g.printfFunctionType(), // Helper method for readability
		g.Ctx.Module.NamedFunction("printf"),
		finalArgs,
		"call.printf",
	)

	return nil
}

func (g *LLVMGenerator) getFormatSpecifier(res *CompilerResult) string {
	var swaType ast.Type
	if res.SwaType != nil {
		swaType = res.SwaType
	} else if res.SymbolTableEntry != nil {
		swaType = res.SymbolTableEntry.DeclaredType
	}

	if swaType != nil {
		switch swaType.(type) {
		case ast.NumberType, *ast.NumberType,
			ast.Number64Type, *ast.Number64Type, *ast.BoolType:
			return "%d"
		case ast.FloatType, *ast.FloatType:
			return "%f"
		case ast.StringType, *ast.StringType:
			return "%s"
		case ast.ErrorType, *ast.ErrorType:
			return "%s"
		}
	}

	// Fallback to LLVM type inspection
	if res.Value != nil {
		lltype := res.Value.Type()
		switch lltype.TypeKind() {
		case llvm.IntegerTypeKind:
			return "%d"
		case llvm.DoubleTypeKind, llvm.FloatTypeKind:
			return "%f"
		case llvm.PointerTypeKind:
			return "%s" // Best guess for pointers in print is string
		}
	}

	return "%p"
}

func (g *LLVMGenerator) printfFunctionType() llvm.Type {
	return llvm.FunctionType(
		llvm.GlobalContext().Int32Type(),
		[]llvm.Type{llvm.PointerType(llvm.GlobalContext().Int8Type(), 0)},
		true, // IsVariadic = true
	)
}

var printableValueExtractors = map[reflect.Type]PrintableValueExtractor{
	// Accessors: These return pointers/addresses and MUST be loaded
	reflect.TypeFor[*ast.MemberExpression]():               extractWithStructType,
	reflect.TypeFor[*ast.ArrayOfStructsAccessExpression](): extractWithStructType,
	reflect.TypeFor[*ast.ArrayAccessExpression]():          extractWithArrayType,

	// Directs: These already hold the value in the result
	reflect.TypeFor[*ast.FunctionCallExpression](): extractDirect,
	reflect.TypeFor[*ast.StringExpression]():       extractDirect,
	reflect.TypeFor[*ast.NumberExpression]():       extractDirect,
	reflect.TypeFor[*ast.FloatExpression]():        extractDirect,
	reflect.TypeFor[*ast.BinaryExpression]():       extractDirect,
	reflect.TypeFor[*ast.SymbolExpression]():       extractSymbol,
	reflect.TypeFor[*ast.PrefixExpression]():       extractDirect,
	reflect.TypeFor[*ast.ZeroExpression]():         extractDirect,
	reflect.TypeFor[*ast.ErrorExpression]():        extractDirect,
	reflect.TypeFor[*ast.BooleanExpression]():      extractDirect,
}

func extractSymbol(g *LLVMGenerator, res *CompilerResult) llvm.Value {
	_, ok := res.SymbolTableEntry.DeclaredType.(ast.ErrorType)

	g.Debugf("IS ErrorType %v", ok)

	if ok {
		return g.Ctx.Builder.CreateLoad(
			res.SymbolTableEntry.Address.Type(),
			*res.SymbolTableEntry.Address,
			"print.array.load",
		)
	}

	return *res.Value
}

func extractDirect(g *LLVMGenerator, res *CompilerResult) llvm.Value {
	val := *res.Value
	// Ensure i1 (bool) is extended for printf
	if val.Type().TypeKind() == llvm.IntegerTypeKind && val.Type().IntTypeWidth() == 1 {
		return g.Ctx.Builder.CreateZExt(val, g.Ctx.Context.Int32Type(), "zext.bool")
	}
	return val
}

func extractWithArrayType(g *LLVMGenerator, res *CompilerResult) llvm.Value {
	return g.Ctx.Builder.CreateLoad(
		res.ArraySymbolTableEntry.UnderlyingType,
		*res.Value,
		"print.array.load",
	)
}
func extractWithStructType(g *LLVMGenerator, res *CompilerResult) llvm.Value {
	return g.Ctx.Builder.CreateLoad(
		*res.StuctPropertyValueType,
		*res.Value,
		"print.struct.load",
	)
}

package compiler

import (
	"fmt"
	"reflect"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

type PrintableValueExtractor func(g *LLVMGenerator, res *ast.CompilerResult) llvm.Value

// VisitPrintStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitPrintStatement(node *ast.PrintStatetement) error {
	printableValues := []llvm.Value{}

	for _, expr := range node.Values {
		err := expr.Accept(g)
		if err != nil {
			return err
		}

		res := g.getLastResult()

		extractor, ok := printableValueExtractors[reflect.TypeOf(expr)]
		if !ok {
			g.NotImplemented(fmt.Sprintf("VisitPrintStatement unimplemented for %T", expr))

			continue
		}

		printableValues = append(printableValues, extractor(g, res))
	}

	// TODO: You must ensure your first argument to printf is the format string!
	// If it's missing from printableValues, printf will likely segfault.
	g.Ctx.Builder.CreateCall(
		g.printfFunctionType(), // Helper method for readability
		g.Ctx.Module.NamedFunction("printf"),
		printableValues,
		"call.printf",
	)

	return nil
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
	reflect.TypeFor[ast.MemberExpression]():               extractWithStructType,
	reflect.TypeFor[ast.ArrayOfStructsAccessExpression](): extractWithStructType,
	reflect.TypeFor[ast.ArrayAccessExpression]():          extractWithArrayType,

	// Directs: These already hold the value in the result
	reflect.TypeFor[ast.FunctionCallExpression](): extractDirect,
	reflect.TypeFor[ast.StringExpression]():       extractDirect,
	reflect.TypeFor[ast.NumberExpression]():       extractDirect,
	reflect.TypeFor[ast.FloatExpression]():        extractDirect,
	reflect.TypeFor[ast.BinaryExpression]():       extractDirect,
	reflect.TypeFor[ast.SymbolExpression]():       extractDirect,
}

func extractDirect(g *LLVMGenerator, res *ast.CompilerResult) llvm.Value {
	return *res.Value
}

func extractWithArrayType(g *LLVMGenerator, res *ast.CompilerResult) llvm.Value {
	return g.Ctx.Builder.CreateLoad(
		res.ArraySymbolTableEntry.UnderlyingType,
		*res.Value,
		"print.array.load",
	)
}
func extractWithStructType(g *LLVMGenerator, res *ast.CompilerResult) llvm.Value {
	return g.Ctx.Builder.CreateLoad(
		*res.StuctPropertyValueType,
		*res.Value,
		"print.struct.load",
	)
}

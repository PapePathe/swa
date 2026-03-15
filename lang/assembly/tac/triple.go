package tac

import (
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

type CustomType struct {
	Types []ast.Type
}

// Triple is the platform-neutral TAC IR.
// GlobalOps holds read-only data (strings) emitted before any procedure.
// Procs holds user-defined functions; Main holds the entry point.
type Triple struct {
	types       map[string]CustomType
	symbolTypes map[string]ast.Type
	GlobalOps   []AsmOp
	lastValue   InstArg
	Main        *Proc
	Procs       []*Proc
	currproc    *Proc
}

var _ AsmOp = (*Triple)(nil)

func (gen *Triple) Gen(g AssemblyOpGenerator) error {
	return g.VisitTriple(gen)
}

func (gen *Triple) Display() {
	fmt.Println("globals:")
	for i, op := range gen.GlobalOps {
		fmt.Printf("  %5d %s\n", i, op)
	}

	fmt.Println()

	for _, proc := range gen.Procs {
		fmt.Printf("%s @%s( ", proc.Ret.String(), proc.Name)
		for _, arg := range proc.Args {
			fmt.Printf("%s:%s ", arg.Name, arg.ArgType.Value())
		}
		fmt.Println(")")

		for _, label := range proc.Labels {
			fmt.Printf("  %s:\n", label.Name)
			for i, op := range label.Ops {
				fmt.Printf("    %5d %+v\n", i, op)
			}
		}
		fmt.Println()
	}

	fmt.Println("main:")
	for _, label := range gen.Main.Labels {
		fmt.Printf("  %s:\n", label.Name)
		for i, op := range label.Ops {
			fmt.Printf("    %5d %+v\n", i, op)
		}
	}
}

var _ ast.CodeGenerator = (*Triple)(nil)

func NewTripleGenerator() *Triple {
	return &Triple{
		GlobalOps:   []AsmOp{},
		types:       map[string]CustomType{},
		symbolTypes: make(map[string]ast.Type),
	}
}

func (gen Triple) getLastValue() InstArg {
	val := gen.lastValue
	gen.lastValue = nil
	return val
}

func (gen *Triple) setLastValue(a InstArg) {
	gen.lastValue = a
}

func appendOpToCurrentProc(gen *Triple, op AsmOp) {
	gen.currproc.AppendOp(op)
}

// opIndex returns an InstArg that identifies the last op appended to the
// current label. Used to wire the result of one op as an operand of another.
func opIndex(gen *Triple) *OpRef {
	idx := gen.currproc.LastOpIndex()
	return &OpRef{idx: idx}
}

func (gen *Triple) VisitByteType(node *ast.ByteType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) ZeroOfByteType(node *ast.ByteType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitArrayType(node *ast.ArrayType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	err := node.Value.Accept(gen)
	if err != nil {
		return err
	}

	incomingValue := gen.getLastValue()

	err = node.Assignee.Accept(gen)
	if err != nil {
		return err
	}

	assignee := gen.getLastValue()
	width := 32
	if node.Assignee.VisitedSwaType().Value() == ast.DataTypeNumber64 {
		width = 64
	}

	appendOpToCurrentProc(gen, &InstWrite{
		Src:   incomingValue,
		Dst:   assignee,
		Width: width,
	})

	return nil
}

func (gen *Triple) VisitBinaryExpression(node *ast.BinaryExpression) error {
	err := node.Left.Accept(gen)
	if err != nil {
		return err
	}

	left := gen.getLastValue()

	err = node.Right.Accept(gen)
	if err != nil {
		return err
	}

	right := gen.getLastValue()

	var op AsmOp
	width := 32

	// Try to infer width from operands if the node type is not set
	leftType := node.Left.VisitedSwaType()
	rightType := node.Right.VisitedSwaType()

	if node.VisitedSwaType() != nil && node.VisitedSwaType().Value() == ast.DataTypeNumber64 {
		width = 64
	} else if (leftType != nil && leftType.Value() == ast.DataTypeNumber64) ||
		(rightType != nil && rightType.Value() == ast.DataTypeNumber64) {
		width = 64
	}

	// Update node type if we inferred it was 64-bit
	if width == 64 && node.SwaType == nil {
		node.SwaType = ast.Number64Type{}
	}

	switch node.Operator.Kind {
	case lexer.Plus:
		op = &InstAdd{Left: left, Right: right, Width: width}
	case lexer.Minus:
		op = &InstSub{Left: left, Right: right, Width: width}
	case lexer.Star:
		op = &InstMul{Left: left, Right: right, Width: width}
	case lexer.Divide:
		op = &InstDiv{Left: left, Right: right, Width: width}
	case lexer.Modulo:
		op = &InstMod{Left: left, Right: right, Width: width}
	default:
		return fmt.Errorf("unsupported binary op %s", node.Operator.Kind)
	}

	appendOpToCurrentProc(gen, op)
	gen.setLastValue(opIndex(gen))

	return nil
}

func (gen *Triple) VisitBlockStatement(node *ast.BlockStatement) error {
	for _, stmt := range node.Body {
		err := stmt.Accept(gen)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gen *Triple) VisitBoolType(node *ast.BoolType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitBooleanExpression(node *ast.BooleanExpression) error {
	if node.Value {
		gen.setLastValue(BoolVal{value: "true"})
		return nil
	}

	gen.setLastValue(BoolVal{value: "false"})

	return nil
}

func (gen *Triple) VisitCallExpression(node *ast.CallExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitErrorExpression(node *ast.ErrorExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitErrorType(node *ast.ErrorType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	return node.Exp.Accept(gen)
}

func (gen *Triple) VisitFloatExpression(node *ast.FloatExpression) error {
	res := Float64Val{value: node.Value}
	gen.setLastValue(&res)

	return nil
}

func (gen *Triple) VisitFloatType(node *ast.FloatType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitFloatingBlockExpression(node *ast.FloatingBlockExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	lab := &Label{Name: "default", Ops: []AsmOp{}}
	proc := &Proc{
		Name:         node.Name,
		Ret:          node.ReturnType,
		Args:         node.Args,
		Labels:       []*Label{lab},
		currentLabel: lab,
		lmap:         map[string]int{"default": 0},
	}
	gen.Procs = append(gen.Procs, proc)
	gen.currproc = proc

	err := node.Body.Accept(gen)
	if err != nil {
		return err
	}

	gen.currproc = nil

	return nil
}

func (gen *Triple) VisitMainStatement(node *ast.MainStatement) error {
	lab := &Label{Name: "default", Ops: []AsmOp{}}
	proc := &Proc{
		Name:         "main",
		Labels:       []*Label{lab},
		currentLabel: lab,
		lmap:         map[string]int{"default": 0},
	}
	gen.Main = proc
	gen.currproc = gen.Main

	err := node.Body.Accept(gen)
	if err != nil {
		return err
	}

	gen.currproc = nil

	return nil
}

func (gen *Triple) VisitMemberExpression(node *ast.MemberExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitNumber64Type(node *ast.Number64Type) error { return nil }

func (gen *Triple) VisitNumberExpression(node *ast.NumberExpression) error {
	// If it's a large number or explicitly typed 64-bit, use Number64Val
	if (node.SwaType != nil && node.SwaType.Value() == ast.DataTypeNumber64) ||
		node.Value > 2147483647 || node.Value < -2147483648 {
		gen.setLastValue(&Number64Val{value: node.Value})
		if node.SwaType == nil {
			node.SwaType = ast.Number64Type{}
		}
	} else {
		// Truncate to int for 32-bit (literals in binary ops/assignments).
		gen.setLastValue(&Number32Val{value: int(node.Value)})
		if node.SwaType == nil {
			node.SwaType = ast.NumberType{}
		}
	}

	return nil
}

func (gen *Triple) VisitNumberType(node *ast.NumberType) error   { return nil }
func (gen *Triple) VisitPointerType(node *ast.PointerType) error { return nil }

func (gen *Triple) VisitPrefixExpression(node *ast.PrefixExpression) error {
	err := node.RightExpression.Accept(gen)
	if err != nil {
		return err
	}

	res := gen.getLastValue()

	switch node.Operator.Kind {
	case lexer.Minus:
		appendOpToCurrentProc(gen, &InstSub{
			Left:  &Number32Val{value: 0},
			Right: res,
			Width: 32,
		})
		gen.setLastValue(opIndex(gen))

		return nil
	default:
		return fmt.Errorf("unsupported op (%s) in PrefixExpression", node.Operator.Kind)
	}
}

func (gen *Triple) VisitPrintStatement(node *ast.PrintStatetement) error {
	// Each value in the print statement maps to a FunCallArg followed by printf.
	for _, expr := range node.Values {
		err := expr.Accept(gen)
		if err != nil {
			return err
		}

		val := gen.getLastValue()
		width := 32
		if expr.VisitedSwaType() != nil && expr.VisitedSwaType().Value() == ast.DataTypeNumber64 {
			width = 64
		}

		appendOpToCurrentProc(gen, &InstFunCallArg{Val: val, Width: width})
	}

	appendOpToCurrentProc(gen, &InstFunCall{Symbol: "printf"})

	return nil
}

func (gen *Triple) VisitReturnStatement(node *ast.ReturnStatement) error {
	err := node.Value.Accept(gen)
	if err != nil {
		return err
	}

	last := gen.getLastValue()

	appendOpToCurrentProc(gen, &Ret{Val: last})

	return nil
}

func (gen *Triple) VisitStringExpression(node *ast.StringExpression) error {
	id := uint32(len(gen.GlobalOps))

	gen.GlobalOps = append(gen.GlobalOps, &InstGlobal{
		ID:    id,
		Value: node.Value,
	})

	gen.setLastValue(&GlobalId{id: id})

	return nil
}

func (gen *Triple) VisitStringType(node *ast.StringType) error { return nil }

func (gen *Triple) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	gen.types[node.Name] = CustomType{Types: node.Types}

	return nil
}

func (gen *Triple) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitSymbolAdressExpression(node *ast.SymbolAdressExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitSymbolExpression(node *ast.SymbolExpression) error {
	gen.setLastValue(&SymbolVal{value: node.Value})

	if t, ok := gen.symbolTypes[node.Value]; ok {
		node.SwaType = t
	}

	return nil
}

func (gen *Triple) VisitSymbolType(node *ast.SymbolType) error { return nil }

func (gen *Triple) VisitSymbolValueExpression(node *ast.SymbolValueExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitTupleAssignmentExpression(node *ast.TupleAssignmentExpression) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitTupleExpression(node *ast.TupleExpression) error {
	return nil
}

func (gen *Triple) VisitTupleType(node *ast.TupleType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	if node.Value == nil {
		node.Value = &ast.ZeroExpression{
			T: node.ExplicitType,
		}
	}

	err := node.Value.Accept(gen)
	if err != nil {
		return err
	}

	value := gen.getLastValue()
	// FIXME size should come from type
	width := 32

	if node.ExplicitType != nil &&
		node.ExplicitType.Value() == ast.DataTypeNumber64 {
		width = 64
	}

	appendOpToCurrentProc(gen, &InstAlloc{
		T:    node.ExplicitType,
		Name: node.Name,
	})

	switch node.ExplicitType.Value() {
	case ast.DataTypeNumber, ast.DataTypeNumber64:
		appendOpToCurrentProc(gen, &InstWrite{
			Src:   value,
			Dst:   &SymbolVal{value: node.Name},
			Width: width,
		})
	case ast.DataTypeFloat:
		appendOpToCurrentProc(gen, &InstWriteFloat{
			Src:   value,
			Dst:   &SymbolVal{value: node.Name},
			Width: width,
		})
	default:
		return fmt.Errorf("Datatype %T not supported", node.ExplicitType)
	}

	if node.ExplicitType != nil {
		gen.symbolTypes[node.Name] = node.ExplicitType
	}

	return nil
}

func (gen *Triple) VisitVoidType(node *ast.VoidType) error { return nil }

func (gen *Triple) VisitWhileStatement(node *ast.WhileStatement) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) VisitZeroExpression(node *ast.ZeroExpression) error {
	return node.T.AcceptZero(gen)
}

func (gen *Triple) ZeroOfArrayType(node *ast.ArrayType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) ZeroOfBoolType(node *ast.BoolType) error {
	gen.setLastValue(BoolVal{value: "false"})

	return nil
}

func (gen *Triple) ZeroOfErrorType(node *ast.ErrorType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) ZeroOfFloatType(node *ast.FloatType) error {
	gen.setLastValue(&Float64Val{value: 0.0})

	return nil
}

func (gen *Triple) ZeroOfNumber64Type(node *ast.Number64Type) error {
	gen.setLastValue(&Number32Val{value: 0})

	return nil
}

func (gen *Triple) ZeroOfNumberType(node *ast.NumberType) error {
	gen.setLastValue(&Number32Val{value: 0})

	return nil
}

func (gen *Triple) ZeroOfPointerType(node *ast.PointerType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) ZeroOfStringType(node *ast.StringType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) ZeroOfSymbolType(node *ast.SymbolType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) ZeroOfTupleType(node *ast.TupleType) error {
	return fmt.Errorf("unimplemented")
}

func (gen *Triple) ZeroOfVoidType(node *ast.VoidType) error {
	return fmt.Errorf("unimplemented")
}

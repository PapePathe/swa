package tac

import (
	"fmt"
	"swahili/lang/ast"
)

type InstArg interface {
	InstructionArg() string
}

type BoolVal struct {
	value string
}

var _ InstArg = (*BoolVal)(nil)

func (b BoolVal) InstructionArg() string {
	return b.value
}

type InstID struct {
	id uint32
}

var _ InstArg = (*InstID)(nil)

func (i InstID) InstructionArg() string {
	return fmt.Sprintf("(%d)", i.id)
}

type TypeID struct {
	T ast.Type
}

var _ InstArg = (*TypeID)(nil)

func (t *TypeID) InstructionArg() string {
	return t.T.String()
}

type Inst struct {
	Operation Op
	ArgOne    InstArg
	ArgTwo    InstArg
}

type Proc struct {
	Name  string
	Ret   ast.Type
	Args  []ast.Type
	Insts []Inst
	Table map[string]InstID
}

func (p Proc) Append(i Inst) {
	p.Insts = append(p.Insts, i)
}

func (p Proc) LastInstID() InstID {
	id := uint32(len(p.Insts) - 1)

	return InstID{id: id}
}

type Triple struct {
	Insts     []Inst
	lastValue InstArg
	Main      *Proc
	Procs     []*Proc
	currproc  *Proc
}

var _ ast.CodeGenerator = (*Triple)(nil)

func NewTripleGenerator() *Triple {
	return &Triple{}
}

func (gen Triple) getLastValue() InstArg {
	val := gen.lastValue
	gen.lastValue = nil

	return val
}

func (gen *Triple) setLastValue(a InstArg) {
	gen.lastValue = a
}

func appendToCurrentProc(gen *Triple, i Inst) {
	gen.currproc.Insts = append(gen.currproc.Insts, i)
}

func getLastInstID(gen *Triple) InstID {
	if len(gen.currproc.Insts) == 0 {
		return InstID{id: uint32(0)}
	}

	return InstID{id: uint32(len(gen.currproc.Insts) - 1)}
}

func (gen *Triple) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	panic("unimplemented")
}

func (gen *Triple) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	panic("unimplemented")
}

// VisitArrayOfStructsAccessExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	panic("unimplemented")
}

// VisitArrayType implements [ast.CodeGenerator].
func (gen *Triple) VisitArrayType(node *ast.ArrayType) error {
	panic("unimplemented")
}

// VisitAssignmentExpression implements [ast.CodeGenerator].
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

	appendToCurrentProc(gen, Inst{
		Operation: OpWrite,
		ArgOne:    incomingValue,
		ArgTwo:    assignee,
	})

	return nil
}

// VisitBinaryExpression implements [ast.CodeGenerator].
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

	op, err := opMap(node.Operator.Kind)
	if err != nil {
		return err
	}

	appendToCurrentProc(gen, Inst{
		Operation: op,
		ArgOne:    left,
		ArgTwo:    right,
	})

	lastInstID := getLastInstID(gen)
	gen.setLastValue(lastInstID)

	return nil
}

// VisitBlockStatement implements [ast.CodeGenerator].
func (gen *Triple) VisitBlockStatement(node *ast.BlockStatement) error {
	for _, stmt := range node.Body {
		err := stmt.Accept(gen)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitBoolType implements [ast.CodeGenerator].
func (gen *Triple) VisitBoolType(node *ast.BoolType) error {
	panic("unimplemented")
}

// VisitBooleanExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitBooleanExpression(node *ast.BooleanExpression) error {
	if node.Value {
		gen.setLastValue(BoolVal{value: "true"})

		return nil
	}

	gen.setLastValue(BoolVal{value: "false"})

	return nil
}

// VisitCallExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitCallExpression(node *ast.CallExpression) error {
	panic("unimplemented")
}

// VisitConditionalStatement implements [ast.CodeGenerator].
func (gen *Triple) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	panic("unimplemented")
}

// VisitErrorExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitErrorExpression(node *ast.ErrorExpression) error {
	panic("unimplemented")
}

// VisitErrorType implements [ast.CodeGenerator].
func (gen *Triple) VisitErrorType(node *ast.ErrorType) error {
	panic("unimplemented")
}

// VisitExpressionStatement implements [ast.CodeGenerator].
func (gen *Triple) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	return node.Exp.Accept(gen)
}

func (gen *Triple) VisitFloatExpression(node *ast.FloatExpression) error {
	gen.setLastValue(node)

	return nil
}

// VisitFloatType implements [ast.CodeGenerator].
func (gen *Triple) VisitFloatType(node *ast.FloatType) error {
	panic("unimplemented")
}

// VisitFloatingBlockExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitFloatingBlockExpression(node *ast.FloatingBlockExpression) error {
	panic("unimplemented")
}

// VisitFunctionCall implements [ast.CodeGenerator].
func (gen *Triple) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	panic("unimplemented")
}

// VisitFunctionDefinition implements [ast.CodeGenerator].
func (gen *Triple) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	proc := &Proc{
		Name:  node.Name,
		Ret:   node.ReturnType,
		Insts: []Inst{},
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
	gen.Main = &Proc{}
	gen.currproc = gen.Main

	err := node.Body.Accept(gen)
	if err != nil {
		return err
	}

	gen.currproc = nil

	return nil
}

// VisitMemberExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitMemberExpression(node *ast.MemberExpression) error {
	panic("unimplemented")
}

// VisitNumber64Type implements [ast.CodeGenerator].
func (gen *Triple) VisitNumber64Type(node *ast.Number64Type) error {
	panic("unimplemented")
}

func (gen *Triple) VisitNumberExpression(node *ast.NumberExpression) error {
	gen.setLastValue(node)

	return nil
}

// VisitNumberType implements [ast.CodeGenerator].
func (gen *Triple) VisitNumberType(node *ast.NumberType) error {
	panic("unimplemented")
}

// VisitPointerType implements [ast.CodeGenerator].
func (gen *Triple) VisitPointerType(node *ast.PointerType) error {
	panic("unimplemented")
}

// VisitPrefixExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitPrefixExpression(node *ast.PrefixExpression) error {
	panic("unimplemented")
}

// VisitPrintStatement implements [ast.CodeGenerator].
func (gen *Triple) VisitPrintStatement(node *ast.PrintStatetement) error {
	panic("unimplemented")
}

// VisitReturnStatement implements [ast.CodeGenerator].
func (gen *Triple) VisitReturnStatement(node *ast.ReturnStatement) error {
	err := node.Value.Accept(gen)
	if err != nil {
		return err
	}

	last := gen.getLastValue()

	appendToCurrentProc(gen, Inst{
		Operation: OpReturn,
		ArgOne:    last,
	})

	return nil
}

// VisitStringExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitStringExpression(node *ast.StringExpression) error {
	panic("unimplemented")
}

// VisitStringType implements [ast.CodeGenerator].
func (gen *Triple) VisitStringType(node *ast.StringType) error {
	panic("unimplemented")
}

// VisitStructDeclaration implements [ast.CodeGenerator].
func (gen *Triple) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	panic("unimplemented")
}

// VisitStructInitializationExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	panic("unimplemented")
}

// VisitSymbolAdressExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitSymbolAdressExpression(node *ast.SymbolAdressExpression) error {
	panic("unimplemented")
}

// VisitSymbolExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitSymbolExpression(node *ast.SymbolExpression) error {
	gen.setLastValue(node)

	return nil
}

// VisitSymbolType implements [ast.CodeGenerator].
func (gen *Triple) VisitSymbolType(node *ast.SymbolType) error {
	panic("unimplemented")
}

// VisitSymbolValueExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitSymbolValueExpression(node *ast.SymbolValueExpression) error {
	panic("unimplemented")
}

// VisitTupleAssignmentExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitTupleAssignmentExpression(node *ast.TupleAssignmentExpression) error {
	panic("unimplemented")
}

// VisitTupleExpression implements [ast.CodeGenerator].
func (gen *Triple) VisitTupleExpression(node *ast.TupleExpression) error {
	panic("unimplemented")
}

// VisitTupleType implements [ast.CodeGenerator].
func (gen *Triple) VisitTupleType(node *ast.TupleType) error {
	panic("unimplemented")
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

	appendToCurrentProc(gen, Inst{
		Operation: OpAlloc,
		ArgOne:    &TypeID{T: node.ExplicitType},
		ArgTwo:    ast.SymbolExpression{Value: node.Name},
	})

	allocInstId := getLastInstID(gen)

	appendToCurrentProc(gen, Inst{
		Operation: OpWrite,
		ArgOne:    value,
		ArgTwo:    allocInstId,
	})

	return nil
}

// VisitVoidType implements [ast.CodeGenerator].
func (gen *Triple) VisitVoidType(node *ast.VoidType) error {
	panic("unimplemented")
}

// VisitWhileStatement implements [ast.CodeGenerator].
func (gen *Triple) VisitWhileStatement(node *ast.WhileStatement) error {
	panic("unimplemented")
}

func (gen *Triple) VisitZeroExpression(node *ast.ZeroExpression) error {
	err := node.T.AcceptZero(gen)
	if err != nil {
		return err
	}

	return nil
}

// ZeroOfArrayType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfArrayType(node *ast.ArrayType) error {
	panic("unimplemented")
}

// ZeroOfBoolType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfBoolType(node *ast.BoolType) error {
	panic("unimplemented")
}

// ZeroOfErrorType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfErrorType(node *ast.ErrorType) error {
	panic("unimplemented")
}

// ZeroOfFloatType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfFloatType(node *ast.FloatType) error {
	gen.setLastValue(ast.FloatExpression{Value: 0.0})

	return nil
}

// ZeroOfNumber64Type implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfNumber64Type(node *ast.Number64Type) error {
	gen.setLastValue(ast.NumberExpression{Value: 0.0})

	return nil
}

// ZeroOfNumberType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfNumberType(node *ast.NumberType) error {
	gen.setLastValue(ast.NumberExpression{Value: 0})

	return nil
}

// ZeroOfPointerType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfPointerType(node *ast.PointerType) error {
	panic("unimplemented")
}

// ZeroOfStringType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfStringType(node *ast.StringType) error {
	panic("unimplemented")
}

// ZeroOfSymbolType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfSymbolType(node *ast.SymbolType) error {
	panic("unimplemented")
}

// ZeroOfTupleType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfTupleType(node *ast.TupleType) error {
	panic("unimplemented")
}

// ZeroOfVoidType implements [ast.CodeGenerator].
func (gen *Triple) ZeroOfVoidType(node *ast.VoidType) error {
	panic("unimplemented")
}

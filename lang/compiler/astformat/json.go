package astformat

import (
	"swahili/lang/ast"
)

type Json struct {
	Element map[string]any
}

func (j *Json) ZeroOfBoolType(node *ast.BoolType) error {
	panic("unimplemented")
}

func (j *Json) VisitBoolType(node *ast.BoolType) error {
	res := make(map[string]any)
	res["BoolType"] = node.Value().String()

	j.setLastResult(res)

	return nil
}

// VisitBooleanExpression implements [ast.CodeGenerator].
func (j *Json) VisitBooleanExpression(node *ast.BooleanExpression) error {
	res := make(map[string]any)
	res["BooleanExpression"] = node.Value

	j.setLastResult(res)

	return nil
}

func NewJsonFormatter() *Json {
	return &Json{Element: map[string]any{}}
}

var _ ast.CodeGenerator = (*Json)(nil)

func (j *Json) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	m := make(map[string]any)

	_ = node.Name.Accept(j)
	m["Name"] = j.getLastResult()

	_ = node.Index.Accept(j)
	m["Index"] = j.getLastResult()

	res := make(map[string]any)
	res["ArrayAccessExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	m := make(map[string]any)

	_ = node.Underlying.Accept(j)
	m["UnderlyingType"] = j.getLastResult()
	m["Contents"] = visitValuesArray(j, node.Contents)

	res := make(map[string]any)
	res["ArrayInitializationExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	m := make(map[string]any)

	_ = node.Name.Accept(j)
	m["Name"] = j.getLastResult()

	_ = node.Index.Accept(j)
	m["Index"] = j.getLastResult()

	_ = node.Property.Accept(j)
	m["Property"] = j.getLastResult()

	res := make(map[string]any)
	res["ArrayOfStructsAccessExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitArrayType(node *ast.ArrayType) error {
	res := make(map[string]any)
	_ = node.Underlying.Accept(j)
	res["ArrayType"] = j.getLastResult()

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	m := make(map[string]any)
	m["Operator"] = node.Operator.Value

	_ = node.Assignee.Accept(j)
	m["Assignee"] = j.getLastResult()

	_ = node.Value.Accept(j)
	m["Value"] = j.getLastResult()

	res := make(map[string]any)
	res["AssignmentExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitTupleExpression(node *ast.TupleExpression) error {
	m := make(map[string]any)
	m["Expressions"] = visitValuesArray(j, node.Expressions)

	res := make(map[string]any)
	res["TupleExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitErrorExpression(node *ast.ErrorExpression) error {
	res := make(map[string]any)
	res["ErrorExpression"] = nil // Error expressions typically don't have a value to format

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitTupleAssignmentExpression(node *ast.TupleAssignmentExpression) error {
	m := make(map[string]any)
	_ = node.Assignees.Accept(j)
	m["Assignees"] = j.getLastResult()
	_ = node.Value.Accept(j)
	m["Value"] = j.getLastResult()

	res := make(map[string]any)
	res["TupleAssignmentExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitBinaryExpression(node *ast.BinaryExpression) error {
	m := make(map[string]any)
	_ = node.Left.Accept(j)
	m["Left"] = j.getLastResult()

	_ = node.Right.Accept(j)
	m["Right"] = j.getLastResult()
	m["Operator"] = node.Operator.Value

	res := make(map[string]any)
	res["BinaryExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitBlockStatement(node *ast.BlockStatement) error {
	res := make(map[string]any)
	res["BlockStatement"] = visitStatementsArray(j, node.Body)

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitCallExpression(node *ast.CallExpression) error {
	m := make(map[string]any)
	m["Caller"] = node.Caller
	m["Arguments"] = node.Arguments

	res := make(map[string]any)
	res["CallExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	m := make(map[string]any)

	_ = node.Success.Accept(j)
	m["success"] = j.getLastResult()

	_ = node.Condition.Accept(j)
	m["condition"] = j.getLastResult()

	_ = node.Failure.Accept(j)
	m["failure"] = j.getLastResult()

	res := make(map[string]any)
	res["ConditionalStatement"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	res := make(map[string]any)

	_ = node.Exp.Accept(j)
	res["ExpressionStatement"] = j.getLastResult()

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitFloatExpression(node *ast.FloatExpression) error {
	res := map[string]any{
		"FloatExpression": node.Value,
	}

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitFloatType(node *ast.FloatType) error {
	res := map[string]any{
		"FloatType": node.Value().String(),
	}

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitFloatingBlockExpression(node *ast.FloatingBlockExpression) error {
	m := make(map[string]any)

	_ = node.Stmt.Accept(j)
	m["Stmt"] = j.getLastResult()

	res := make(map[string]any)
	res["FloatingBlockExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	m := make(map[string]any)
	m["Name"] = node.Name

	m["Args"] = visitValuesArray(j, node.Args)

	res := make(map[string]any)
	res["FunctionCallExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	m := make(map[string]any)
	m["Name"] = node.Name

	_ = node.ReturnType.Accept(j)
	m["ReturnType"] = j.getLastResult()
	m["Args"] = node.Args

	_ = node.Body.Accept(j)
	m["Body"] = j.getLastResult()

	res := make(map[string]any)
	res["FunctionDeclarationStatement"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitMainStatement(node *ast.MainStatement) error {
	m := make(map[string]any)
	_ = node.Body.Accept(j)
	m["Body"] = j.getLastResult()

	res := make(map[string]any)
	res["MainProGram"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitMemberExpression(node *ast.MemberExpression) error {
	m := make(map[string]any)

	_ = node.Object.Accept(j)
	m["Object"] = j.getLastResult()

	_ = node.Property.Accept(j)
	m["Property"] = j.getLastResult()
	m["Computed"] = node.Computed

	res := make(map[string]any)
	res["MemberExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitNumber64Type(node *ast.Number64Type) error {
	res := make(map[string]any)
	res["Number64Type"] = node.Value()

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitNumberExpression(node *ast.NumberExpression) error {
	res := map[string]any{
		"NumberExpression": node.Value,
	}

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitNumberType(node *ast.NumberType) error {
	res := make(map[string]any)
	res["NumberType"] = node.Value().String()

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitPointerType(node *ast.PointerType) error {
	res := make(map[string]any)
	res["PointerType"] = node.Underlying.Value().String()

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitPrefixExpression(node *ast.PrefixExpression) error {
	m := make(map[string]any)
	m["Operator"] = node.Operator

	_ = node.RightExpression.Accept(j)
	m["RightExpression"] = j.getLastResult()

	res := make(map[string]any)
	res["PrefixExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitPrintStatement(node *ast.PrintStatetement) error {
	res := make(map[string]any)
	res["PrintStatetement"] = visitValuesArray(j, node.Values)

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitReturnStatement(node *ast.ReturnStatement) error {
	res := make(map[string]any)
	res["ReturnStatement"] = node.Value

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitStringExpression(node *ast.StringExpression) error {
	res := map[string]any{
		"StringExpression": node.Value,
	}

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitStringType(node *ast.StringType) error {
	res := map[string]any{
		"StringType": node.Value().String(),
	}

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	m := make(map[string]any)
	m["Name"] = node.Name
	m["Properties"] = node.Properties
	values := []map[string]any{}

	for _, v := range node.Types {
		_ = v.Accept(j)
		values = append(values, j.getLastResult())
	}

	m["Types"] = values

	res := make(map[string]any)
	res["StructDeclarationStatement"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	m := make(map[string]any)
	m["Name"] = node.Name
	m["Properties"] = node.Properties
	m["Values"] = visitValuesArray(j, node.Values)
	res := make(map[string]any)
	res["StructInitializationExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitSymbolExpression(node *ast.SymbolExpression) error {
	res := make(map[string]any)
	res["SymbolExpression"] = node.Value

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitSymbolType(node *ast.SymbolType) error {
	res := make(map[string]any)
	res["SymbolType"] = node.Name

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitSymbolAdressExpression(node *ast.SymbolAdressExpression) error {
	m := make(map[string]any)

	_ = node.Exp.Accept(j)
	m["Exp"] = j.getLastResult()

	res := make(map[string]any)
	res["SymbolAdressExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitSymbolValueExpression(node *ast.SymbolValueExpression) error {
	m := make(map[string]any)

	_ = node.Exp.Accept(j)
	m["Exp"] = j.getLastResult()

	res := make(map[string]any)
	res["SymbolValueExpression"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitErrorType(node *ast.ErrorType) error {
	res := make(map[string]any)
	res["ErrorType"] = node.Value().String()
	j.setLastResult(res)
	return nil
}

func (j *Json) VisitTupleType(node *ast.TupleType) error {
	m := make(map[string]any)
	types := []map[string]any{}
	for _, t := range node.Types {
		_ = t.Accept(j)
		types = append(types, j.getLastResult())
	}
	m["Types"] = types
	res := make(map[string]any)
	res["TupleType"] = m
	j.setLastResult(res)
	return nil
}

func (j *Json) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	m := make(map[string]any)

	m["Name"] = node.Name
	m["IsConstant"] = node.IsConstant

	if node.Value != nil {
		_ = node.Value.Accept(j)
		m["Value"] = j.getLastResult()
	}

	_ = node.ExplicitType.Accept(j)
	m["ExplicitType"] = j.getLastResult()

	res := make(map[string]any)
	res["VarDeclarationStatement"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitVoidType(node *ast.VoidType) error {
	res := make(map[string]any)
	res["VoidType"] = node.Value().String()

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitWhileStatement(node *ast.WhileStatement) error {
	m := make(map[string]any)

	_ = node.Condition.Accept(j)
	m["Condition"] = j.getLastResult()

	_ = node.Body.Accept(j)
	m["Body"] = j.getLastResult()

	res := make(map[string]any)
	res["WhileStatement"] = m

	j.setLastResult(res)

	return nil
}

func (j *Json) VisitZeroExpression(node *ast.ZeroExpression) error {
	res := make(map[string]any)
	res["ZeroExpression"] = node.T.Value().String()

	j.setLastResult(res)

	return nil
}

func (j *Json) ZeroOfArrayType(node *ast.ArrayType) error {
	panic("unimplemented")
}

func (j *Json) ZeroOfFloatType(node *ast.FloatType) error {
	panic("unimplemented")
}

func (j *Json) ZeroOfNumber64Type(node *ast.Number64Type) error {
	panic("unimplemented")
}

func (j *Json) ZeroOfNumberType(node *ast.NumberType) error {
	panic("unimplemented")
}

func (j *Json) ZeroOfPointerType(node *ast.PointerType) error {
	return nil
}

func (j *Json) ZeroOfErrorType(node *ast.ErrorType) error {
	return nil
}

func (j *Json) ZeroOfTupleType(node *ast.TupleType) error {
	return nil
}

func (j *Json) ZeroOfStringType(node *ast.StringType) error {
	panic("unimplemented")
}

func (j *Json) ZeroOfSymbolType(node *ast.SymbolType) error {
	panic("unimplemented")
}

func (j *Json) ZeroOfVoidType(node *ast.VoidType) error {
	panic("unimplemented")
}

func (g *Json) setLastResult(res map[string]any) {
	g.Element = res
}

func (g *Json) getLastResult() map[string]any {
	r := g.Element

	g.Element = nil

	return r
}

func visitValuesArray(j *Json, arr []ast.Expression) []map[string]any {
	values := []map[string]any{}

	for _, v := range arr {
		_ = v.Accept(j)
		values = append(values, j.getLastResult())
	}

	return values
}

func visitStatementsArray(j *Json, arr []ast.Statement) []map[string]any {
	values := []map[string]any{}

	for _, v := range arr {
		_ = v.Accept(j)
		values = append(values, j.getLastResult())
	}

	return values
}

package tac

import (
	"fmt"
	"swahili/lang/lexer"
)

type Op int

func (op Op) String() string {
	switch op {
	case OpModulo:
		return "modulo"
	case OpMultiply:
		return "muliply"
	case OpMinus:
		return "subtract"
	case OpAdd:
		return "add"
	case OpWrite:
		return "write"
	case OpReturn:
		return "return"
	case OpDivide:
		return "divide"
	case OpAlloc:
		return "allocate"
	case OpRead:
		return "read"
	case OpGreaterThanEq:
		return "cmp-greaterthaneq"
	case OpGreaterThan:
		return "cmp-greaterthan"
	case OpLessThanEQ:
		return "cmp-lessthaneq"
	case OpLessThan:
		return "cmp-lessthan"
	case OpEquals:
		return "cmp-equals"
	case OpFunCall:
		return "function-call"
	case OpFunCallArg:
		return "function-call-arg"
	case OpGlobal:
		return "global"
	case OpNegation:
		return "negate"
	case OpJumpCond:
		return "conditional-jump"
	case OpJump:
		return "jump"
	default:
		return "unknown op"
	}
}

const (
	OpAdd = iota
	OpAlloc
	OpDivide
	OpEquals
	OpGlobal
	OpMinus
	OpGreaterThan
	OpGreaterThanEq
	OpLessThan
	OpLessThanEQ
	OpMultiply
	OpModulo
	OpNegation
	OpFunCallArg
	OpFunCall
	OpWrite
	OpRead
	OpReturn
	OpJump
	OpJumpCond
)

func opMap(tk lexer.TokenKind) (Op, error) {
	switch tk {
	case lexer.Star:
		return OpMultiply, nil
	case lexer.Plus:
		return OpAdd, nil
	case lexer.Minus:
		return OpMinus, nil
	case lexer.Divide:
		return OpDivide, nil
	case lexer.Modulo:
		return OpModulo, nil
	case lexer.Equals:
		return OpEquals, nil
	case lexer.GreaterThanEquals:
		return OpGreaterThanEq, nil
	case lexer.GreaterThan:
		return OpGreaterThan, nil
	case lexer.LessThanEquals:
		return OpLessThanEQ, nil
	case lexer.LessThan:
		return OpLessThan, nil
	default:
		return 0, fmt.Errorf("unsupported op map %s", tk)
	}
}

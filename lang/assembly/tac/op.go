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
	case OpEquals:
		return "equals"
	default:
		return "unknown op"
	}
}

const (
	OpAdd = iota
	OpAlloc
	OpDivide
	OpEquals
	OpMinus
	OpMultiply
	OpModulo
	OpWrite
	OpRead
	OpReturn
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
	default:
		return 0, fmt.Errorf("unsupported op map %s", tk)
	}
}

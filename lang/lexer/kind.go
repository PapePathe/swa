package lexer

// TokenKind is a custom type representing the kind of a token.
type TokenKind int

// String returns a string representation of the TokenKind.
func (k TokenKind) String() string {
	switch k {
	case Struct:
		return "STRUCT"
	case Let:
		return "LET"
	case Const:
		return "CONST"
	case Or:
		return "OR"
	case And:
		return "AND"
	case KeywordIf:
		return "IF"
	case KeywordElse:
		return "ELSE"
	case TypeInt:
		return "INT"
	case Identifier:
		return "IDENTIFIER"
	case Assignment:
		return "ASSIGNMENT"
	case String:
		return "STRING"
	case Colon:
		return "COLON"
	case Comma:
		return "COMMA"
	case CloseCurly:
		return "CLOSE_CURLY"
	case CloseParen:
		return "CLOSE_PAREN"
	case CloseBracket:
		return "CLOSE_BRACKET"
	case Divide:
		return "DIVIDE"
	case EOF:
		return "EOF"
	case Equals:
		return "EQUALS"
	case GreaterThan:
		return "GREATER_THAN"
	case GreaterThanEquals:
		return "GREATER_THAN_EQUALS"
	case LessThan:
		return "LESS_THAN"
	case LessThanEquals:
		return "LESS_THAN_EQUALS"
	case Minus:
		return "MINUS"
	case Multiply:
		return "MULTIPLY"
	case Not:
		return "NOT"
	case NotEquals:
		return "NOT_EQUALS"
	case OpenCurly:
		return "OPEN_CURLY"
	case OpenParen:
		return "OPEN_PAREN"
	case OpenBracket:
		return "OPEN_BRACKET"
	case Plus:
		return "PLUS"
	case PlusEquals:
		return "PLUS_EQUAL"
	case Star:
		return "STAR"
	case SemiColon:
		return "SEMI_COLON"
	case QuestionMark:
		return "QUESTION_MARK"
	case Number:
		return "NUMBER"
	}

	return "UNKNOWN TOKEN KIND"
}

const (
	// Let represents the "let" keyword.
	Let TokenKind = iota
	// Const represents the "const" keyword.
	Const
	// Or represents the "or" logical operator.
	Or
	// And represents the "and" logical operator.
	And
	// KeywordIf represents the "if" keyword.
	KeywordIf
	// KeywordElse represents the "else" keyword.
	KeywordElse

	// TypeInt represents the "int" type.
	TypeInt
	// String represents a string literal.
	String
	// Number represents a numeric literal.
	Number

	// Assignment represents an assignment operator.
	Assignment
	// Colon represents a colon symbol.
	Colon
	// Comma represents a comma symbol.
	Comma
	// CloseCurly represents a closing curly brace symbol.
	CloseCurly
	// CloseParen represents a closing parenthesis symbol.
	CloseParen
	// CloseBracket represents a closing bracket symbol.
	CloseBracket
	// Divide represents the divide operator.
	Divide
	// EOF represents the end of file token.
	EOF
	// Equals represents the equality operator.
	Equals
	// GreaterThan represents the greater than operator.
	GreaterThan
	// GreaterThanEquals represents the greater than or equal to operator.
	GreaterThanEquals
	// LessThan represents the less than operator.
	LessThan
	// LessThanEquals represents the less than or equal to operator.
	LessThanEquals
	// Minus represents the minus operator.
	Minus
	// MinusEquals.
	MinusEquals
	// Multiply represents the multiplication operator.
	Multiply
	// Not represents the negation operator.
	Not
	// NotEquals represents the not equal operator.
	NotEquals
	// OpenCurly represents an opening curly brace symbol.
	OpenCurly
	// OpenParen represents an opening parenthesis symbol.
	OpenParen
	// OpenBracket represents an opening bracket symbol.
	OpenBracket
	// Plus represents the addition operator.
	Plus
	// Plus represents the plus equals operator.
	PlusEquals
	// Star represents the multiplication symbol (star).
	Star
	// SemiColon represents a semicolon symbol.
	SemiColon
	// QuestionMark represents a question mark symbol.
	QuestionMark
	// Identifier represents an identifier (variable or function name).
	Identifier
	// Struct ...
	Struct
)

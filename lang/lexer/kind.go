package lexer

// TokenKind is a custom type representing the kind of a token.
type TokenKind int

var tks = map[TokenKind]string{
	And:                "AND",
	Assignment:         "ASSIGNMENT",
	Colon:              "COLON",
	Comma:              "COMMA",
	CloseCurly:         "CLOSE_CURLY",
	CloseParen:         "CLOSE_PAREN",
	CloseBracket:       "CLOSE_BRACKET",
	Const:              "CONST",
	DialectDeclaration: "DIALECT",
	Divide:             "DIVIDE",
	Dot:                "DOT",
	EOF:                "EOF",
	Equals:             "EQUALS",
	Function:           "FUNCTION",
	GreaterThan:        "GREATER_THAN",
	GreaterThanEquals:  "GREATER_THAN_EQUALS",
	Identifier:         "IDENTIFIER",
	KeywordIf:          "IF",
	KeywordElse:        "ELSE",
	KeywordWhile:       "WHILE",
	LessThan:           "LESS_THAN",
	LessThanEquals:     "LESS_THAN_EQUALS",
	Let:                "LET",
	Main:               "MAIN_PROGRAM",
	Minus:              "MINUS",
	Multiply:           "MULTIPLY",
	Modulo:             "MODULO",
	Not:                "NOT",
	NotEquals:          "NOT_EQUALS",
	Number:             "NUMBER",
	Float:              "DECIMAL",
	OpenCurly:          "OPEN_CURLY",
	OpenParen:          "OPEN_PAREN",
	OpenBracket:        "OPEN_BRACKET",
	Or:                 "OR",
	Plus:               "PLUS",
	PlusEquals:         "PLUS_EQUAL",
	Print:              "PRINT",
	QuestionMark:       "QUESTION_MARK",
	SemiColon:          "SEMI_COLON",
	Return:             "RETURN",
	StarEquals:         "STAR_EQUALS",
	Struct:             "STRUCT",
	String:             "STRING",
	Star:               "STAR",
	TypeInt:            "TYPE_INT",
	TypeInt64:          "TYPE_INT_64",
	TypeFloat:          "TYPE_FLOAT",
	TypeString:         "TYPE_STRING",
	TypeError:          "TYPE_ERROR",
	Variadic:           "VARIADIC",
	Zero:               "ZERO_OF",
}

// String s a string representation of the TokenKind.
func (k TokenKind) String() string {
	val, ok := tks[k]
	if !ok {
		return ""
	}

	return val
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
	// KeywordWhile represents the while statement
	KeywordWhile
	// TypeInt represents the "int" type.
	TypeInt
	TypeInt64
	TypeFloat
	// TypeString represents the "string" type.
	TypeString
	// String represents a string literal.
	String
	// Number represents a numeric literal.
	Number
	// Float represents a decimal literal.
	Float
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
	// PlusEquals represents the plus equals operator.
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
	// Dot ...
	Dot
	// StarEquals ...
	StarEquals
	// DialectDeclaration ...
	DialectDeclaration
	// Print
	Print
	// Character
	Character
	Function
	Main
	Return
	Modulo

	Variadic
	TypeError
	Zero
)

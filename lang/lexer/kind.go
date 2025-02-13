package lexer

type TokenKind int

func (k TokenKind) String() string {
	switch k {
	case KEYWORD_IF:
		return "IF"
	case KEYWORD_ELSE:
		return "ELSE"
	case TYPE_INT:
		return "INT"
	case IDENTIFIER:
		return "IDENTIFIER"
	case ASSIGNMENT:
		return "ASSIGNMENT"
	case STRING:
		return "STRING"
	case COLON:
		return "COLON"
	case COMMA:
		return "COMMA"
	case CLOSE_CURLY:
		return "CLOSE_CURLY"
	case CLOSE_PAREN:
		return "CLOSE_PAREN"
	case CLOSE_BRACKET:
		return "CLOSE_BRACKET"
	case DIVIDE:
		return "DIVIDE"
	case EOF:
		return "EOF"
	case EQUALS:
		return "EQUALS"
	case GREATER_THAN:
		return "GREATER_THAN"
	case GREATER_THAN_EQUALS:
		return "GREATER_THAN_EQUALS"
	case LESS_THAN:
		return "LESS_THAN"
	case LESS_THAN_EQUALS:
		return "LESS_THAN_EQUALS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case NOT:
		return "NOT"
	case NOT_EQUALS:
		return "NOT_EQUALS"
	case OPEN_CURLY:
		return "OPEN_CURLY"
	case OPEN_PAREN:
		return "OPEN_PAREN"
	case OPEN_BRACKET:
		return "OPEN_BRACKET"
	case PLUS:
		return "PLUS"
	case STAR:
		return "STAR"
	case SEMI_COLON:
		return "SEMI_COLON"
	case QUESTION_MARK:
		return "QUESTION_MARK"
	case NUMBER:
		return "NUMBER"

	}

	return "UNKNOWN TOKEN KIND"
}

const (
	ASSIGNMENT TokenKind = iota
	COLON
	COMMA
	CLOSE_CURLY
	CLOSE_PAREN
	CLOSE_BRACKET
	DIVIDE
	EOF
	EQUALS
	GREATER_THAN
	GREATER_THAN_EQUALS
	IDENTIFIER
	LESS_THAN
	LESS_THAN_EQUALS
	MINUS
	MULTIPLY
	NOT
	NOT_EQUALS
	OPEN_CURLY
	OPEN_PAREN
	OPEN_BRACKET
	PLUS
	STAR
	SEMI_COLON
	QUESTION_MARK
	NUMBER
	STRING
	TYPE_INT
	KEYWORD_IF
	KEYWORD_ELSE
)

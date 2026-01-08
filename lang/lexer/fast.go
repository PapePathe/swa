package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// FastLexer is a high-performance lexer that doesn't use regular expressions
type FastLexer struct {
	Tokens        []Token              // The tokens
	source        string               // The source code
	position      int                  // The current position of the lexer
	line          int                  // Current line number
	column        int                  // Current column number
	reservedWords map[string]TokenKind // list of reserved words
	dialect       Dialect
	start         int // start position of current token
	startLine     int // line at start of current token
	startColumn   int // column at start of current token
}

// NewFastLexer creates a new FastLexer instance
func NewFastLexer(source string) (*FastLexer, Dialect, error) {
	dialect, err := getDialect(source)
	if err != nil {
		return nil, nil, err
	}

	return &FastLexer{
		Tokens:        make([]Token, 0),
		reservedWords: dialect.Reserved(),
		source:        source,
		dialect:       dialect,
		line:          1,
		column:        1,
	}, dialect, nil
}

// NewFastLexer creates a new FastLexer instance with a specific dialect
func NewFastLexerWithDialect(source string, dialect Dialect) (*FastLexer, error) {
	return &FastLexer{
		Tokens:        make([]Token, 0),
		reservedWords: dialect.Reserved(),
		source:        source,
		dialect:       dialect,
		line:          1,
		column:        1,
	}, nil
}

// GetAllTokens returns all tokens from lexing the source
func (lex *FastLexer) GetAllTokens() ([]Token, error) {
	err := lex.Lex()
	if err != nil {
		return nil, err
	}

	return lex.Tokens, nil
}

// Lex scans the source code and produces tokens
func (lex *FastLexer) Lex() error {
	for !lex.atEOF() {
		lex.start = lex.position
		lex.startLine = lex.line
		lex.startColumn = lex.column
		ch := lex.advance()

		switch {
		case ch == '\n':
			//			lex.push(Token{Kind: Newline, Value: "\n", Line: lex.startLine, Column: lex.startColumn})

		case unicode.IsSpace(ch) && ch != '\n':
			lex.skipWhitespace()

		case lex.isIdentifierStart(ch):
			lex.lexIdentifierOrKeyword()

		case unicode.IsDigit(ch) || (ch == '-' && unicode.IsDigit(lex.peek())):
			lex.lexNumber()

		case ch == '"':
			if err := lex.lexString(); err != nil {
				return err
			}

		case ch == '\'':
			if err := lex.lexChar(); err != nil {
				return err
			}

		case ch == '/':
			if lex.peek() == '/' {
				lex.lexComment()
			} else {
				lex.lexOperator()
			}

		case ch == ';':
			lex.push2(SemiColon, string(ch))

		case ch == '(':
			lex.push2(OpenParen, "(")

		case ch == ')':
			lex.push2(CloseParen, ")")

		case ch == '{':
			lex.push2(OpenCurly, "{")

		case ch == '}':
			lex.push2(CloseCurly, "}")

		case ch == '[':
			lex.push2(OpenBracket, "[")

		case ch == ']':
			lex.push2(CloseBracket, "]")

		default:
			lex.lexOperator()
		}
	}

	return nil
}

func (lex *FastLexer) skipWhitespace() {
	for !lex.atEOF() {
		ch := lex.peek()
		if !unicode.IsSpace(ch) || ch == '\n' {
			break
		}

		lex.advance()
	}
}

func (lex *FastLexer) lexIdentifierOrKeyword() {
	for !lex.atEOF() {
		ch := lex.peek()

		if ch >= utf8.RuneSelf && unicode.IsLetter(ch) {
			lex.advance()

			continue
		}

		if !lex.isIdentifierPart(ch) {
			break
		}

		lex.advance()
	}

	text := lex.currentText()

	if kind, ok := lex.reservedWords[text]; ok {
		lex.push2(kind, text)

		return
	}

	lex.push2(Identifier, text)
}

func (lex *FastLexer) isIdentifierStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func (lex *FastLexer) isIdentifierPart(ch rune) bool {
	return unicode.IsLetter(ch) ||
		unicode.IsDigit(ch) ||
		ch == '_'
}

func (lex *FastLexer) lexNumber() {
	hadDecimal := false

	// Handle negative numbers
	if lex.currentText() == "-" {
		lex.advance() // Skip the minus that started this token
	}

	for !lex.atEOF() {
		ch := lex.peek()
		if ch == '.' && !hadDecimal && unicode.IsDigit(lex.peekNext()) {
			hadDecimal = true
			lex.advance()
		} else if unicode.IsDigit(ch) {
			lex.advance()
		} else {
			break
		}
	}

	text := lex.currentText()
	if hadDecimal {
		lex.push2(Float, text)
	} else {
		lex.push2(Number, text)
	}
}

func (lex *FastLexer) lexString() error {
	var builder strings.Builder
	escaped := false

	for !lex.atEOF() {
		ch := lex.advance()

		if escaped {
			switch ch {
			case 'n':
				builder.WriteRune('\n')
			case 't':
				builder.WriteRune('\t')
			case '"':
				builder.WriteRune('"')
			case '\\':
				builder.WriteRune('\\')
			case 'r':
				builder.WriteRune('\r')
			default:
				// Invalid escape sequence
				return fmt.Errorf("invalid escape sequence \\%c at line %d, column %d", ch, lex.line, lex.column)
			}
			escaped = false
		} else if ch == '\\' {
			escaped = true
		} else if ch == '"' {
			lex.push2(String, builder.String())
			return nil
		} else {
			builder.WriteRune(ch)
		}
	}

	return fmt.Errorf("unterminated string literal at line %d, column %d", lex.startLine, lex.startColumn)
}

func (lex *FastLexer) lexChar() error {
	if lex.atEOF() {
		return fmt.Errorf("unterminated character literal at line %d, column %d", lex.startLine, lex.startColumn)
	}

	ch := lex.advance()
	var charValue rune

	if ch == '\\' {
		// Handle escape sequences
		if lex.atEOF() {
			return fmt.Errorf("unterminated character literal at line %d, column %d", lex.startLine, lex.startColumn)
		}
		next := lex.advance()
		switch next {
		case 'n':
			charValue = '\n'
		case 't':
			charValue = '\t'
		case '\'':
			charValue = '\''
		case '\\':
			charValue = '\\'
		case 'r':
			charValue = '\r'
		default:
			return fmt.Errorf("invalid escape sequence in character literal at line %d, column %d", lex.line, lex.column)
		}
	} else {
		charValue = ch
	}

	// Expect closing quote
	if lex.atEOF() || lex.advance() != '\'' {
		return fmt.Errorf("unterminated character literal at line %d, column %d", lex.startLine, lex.startColumn)
	}

	lex.push2(Character, string(charValue))
	return nil
}

func (lex *FastLexer) lexComment() {
	// Skip the second slash
	lex.advance()

	for !lex.atEOF() {
		ch := lex.peek()
		if ch == '\n' {
			break
		}

		lex.advance()
	}
}

func (lex *FastLexer) lexOperator() {
	text := lex.currentText()

	if !lex.atEOF() {
		next := lex.peek()
		twoChar := text + string(next)

		if lex.isTwoCharOperator(twoChar) {
			lex.advance()
			lex.push2(lex.getOperatorKind(twoChar), twoChar)
		} else {
			lex.push2(lex.getOperatorKind(text), text)
		}
	}
}

func (lex *FastLexer) isTwoCharOperator(op string) bool {
	switch op {
	case "!=", "==", "<=", ">=", "+=", "||", "&&":
		return true
	default:
		return false
	}
}

func (lex *FastLexer) getOperatorKind(op string) TokenKind {
	switch op {
	case "+":
		return Plus
	case "-":
		return Minus
	case "*":
		return Star
	case "/":
		return Divide
	case "%":
		return Modulo
	case "=":
		return Assignment
	case "==":
		return Equals
	case "!=":
		return NotEquals
	case "!":
		return Not
	case "<":
		return LessThan
	case ">":
		return GreaterThan
	case "<=":
		return LessThanEquals
	case ">=":
		return GreaterThanEquals
	case "&&":
		return And
	case "||":
		return Or
	case "+=":
		return PlusEquals
	case ".":
		return Dot
	case ",":
		return Comma
	case ":":
		return Colon

	default:
		panic(fmt.Sprintf("lex not supported (%s) %v", string(op), lex.Tokens))
	}
}

func (lex *FastLexer) advance() rune {
	if lex.position >= len(lex.source) {
		return 0
	}
	ch := rune(lex.source[lex.position])
	lex.position++

	if ch == '\n' {
		lex.line++
		lex.column = 1
	} else {
		lex.column++
	}

	return ch
}

func (lex *FastLexer) peek() rune {
	if lex.position >= len(lex.source) {
		return 0
	}
	return rune(lex.source[lex.position])
}

func (lex *FastLexer) peekNext() rune {
	if lex.position+1 >= len(lex.source) {
		return 0
	}
	return rune(lex.source[lex.position+1])
}

func (lex *FastLexer) newLine() {
	lex.line += 1
	lex.column = 1
}

func (lex *FastLexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *FastLexer) push2(k TokenKind, v string) {
	lex.Tokens = append(lex.Tokens, Token{
		Kind:   k,
		Value:  v,
		Line:   lex.startLine,
		Column: lex.startColumn,
	})
}

func (lex *FastLexer) atEOF() bool {
	return lex.position >= len(lex.source)
}

func (lex *FastLexer) currentText() string {
	return lex.source[lex.start:lex.position]
}

package lexer

import (
	"fmt"
	"regexp"
)

func Tokenize(source string) []Token {
	lex := createLexer(source)

	for !lex.at_eof() {
		matched := false

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())

			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break

			}
		}

		if !matched {
			panic(fmt.Sprintf("Lexer::Error -> unrecognized token near %s", lex.remainder()))
		}
	}

	lex.push(NewToken(EOF, "EOF"))

	return lex.Tokens
}

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *Lexer, _ *regexp.Regexp) {
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value))
	}
}

func commentHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	if match != nil {
		// Advance past the entire comment.
		lex.advanceN(match[1])
		//		lex.line++
	}
}
func stringHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]

	lex.push(NewToken(STRING, stringLiteral))
	lex.advanceN(len(stringLiteral))
}

func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(NUMBER, match))
	lex.advanceN(len(match))
}

func symbolHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if kind, found := reserved_lu[match]; found {
		lex.push(NewToken(kind, match))
	} else {
		lex.push(NewToken(IDENTIFIER, match))
	}

	lex.advanceN(len(match))
}

func skipHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

type RegexpPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type Lexer struct {
	Tokens   []Token
	source   string
	position int
	patterns []RegexpPattern
}

func (lex *Lexer) advanceN(n int) {
	lex.position += n
}

func (lex *Lexer) at() byte {
	return lex.source[lex.position]
}

func (lex *Lexer) advance() {
	lex.position += 1
}

func (lex *Lexer) remainder() string {
	return lex.source[lex.position:]
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) at_eof() bool {
	return lex.position >= len(lex.source)
}

func createLexer(source string) *Lexer {
	return &Lexer{
		position: 0,
		source:   source,

		Tokens: make([]Token, 0),
		patterns: []RegexpPattern{
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`fèndo`), defaultHandler(TYPE_INT, "fèndo")},
			{regexp.MustCompile(`ni`), defaultHandler(KEYWORD_IF, "ni")},
			{regexp.MustCompile(`nii`), defaultHandler(KEYWORD_ELSE, "nii")},
			{regexp.MustCompile(`\/\/.*`), commentHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			//			{regexp.MustCompile(`==`), defaultHandler(PERCENT, "%")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_THAN_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS_THAN, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_THAN_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER_THAN, ">")},
			//			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			//			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			//			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			//			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			//			{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
			{regexp.MustCompile(`/`), defaultHandler(DIVIDE, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			//			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	// "true":    TRUE,
	// "false":   FALSE,
	// "null":    NULL,
	// "let":     LET,
	// "const":   CONST,
	// "class":   CLASS,
	// "new":     NEW,
	// "import":  IMPORT,
	// "from":    FROM,
	// "fn":      FN,
	// "if":      IF,
	// "else":    ELSE,
	// "foreach": FOREACH,
	// "while":   WHILE,
	// "for":     FOR,
	// "export":  EXPORT,
	// "typeof":  TYPEOF,
	// "in":      IN,
}

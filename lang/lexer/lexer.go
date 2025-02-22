package lexer

import (
	"fmt"
	"regexp"
)

// Tokenize ...
func Tokenize(source string) []Token {
	lex := createLexer(source)

	for !lex.atEOF() {
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
		lex.advanceN(match[1])
	}
}

func stringHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]

	lex.push(NewToken(String, stringLiteral))
	lex.advanceN(len(stringLiteral))
}

func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(Number, match))
	lex.advanceN(len(match))
}

func symbolHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if kind, found := reservedLu[match]; found {
		lex.push(NewToken(kind, match))
	} else {
		lex.push(NewToken(Identifier, match))
	}

	lex.advanceN(len(match))
}

func skipHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

// RegexpPattern ...
type RegexpPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

// Lexer ...
type Lexer struct {
	Tokens   []Token         // The tokens
	source   string          // The source code
	position int             // The current position of the lexer
	patterns []RegexpPattern // the list of patterns of the language
}

func (lex *Lexer) advanceN(n int) {
	lex.position += n
}

//func (lex *Lexer) at() byte {
//	return lex.source[lex.position]
//}

//func (lex *Lexer) advance() {
//	lex.position += 1
//}

func (lex *Lexer) remainder() string {
	return lex.source[lex.position:]
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) atEOF() bool {
	return lex.position >= len(lex.source)
}

func createLexer(source string) *Lexer {
	return &Lexer{
		position: 0,
		source:   source,

		Tokens: make([]Token, 0),
		patterns: []RegexpPattern{
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`fèndo`), defaultHandler(TypeInt, "fèndo")},
			{regexp.MustCompile(`nii`), defaultHandler(KeywordElse, "nii")},
			{regexp.MustCompile(`ni`), defaultHandler(KeywordIf, "ni")},
			{regexp.MustCompile(`struct`), defaultHandler(Struct, "struct")},
			{regexp.MustCompile(`\/\/.*`), commentHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
			{regexp.MustCompile(`\[`), defaultHandler(OpenBracket, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CloseBracket, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OpenCurly, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CloseCurly, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OpenParen, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CloseParen, ")")},
			//			{regexp.MustCompile(`==`), defaultHandler(PERCENT, "%")},
			{regexp.MustCompile(`!=`), defaultHandler(NotEquals, "!=")},
			{regexp.MustCompile(`\+=`), defaultHandler(PlusEquals, "+=")},
			{regexp.MustCompile(`=`), defaultHandler(Assignment, "=")},
			{regexp.MustCompile(`!`), defaultHandler(Not, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LessThanEquals, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LessThan, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GreaterThanEquals, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GreaterThan, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(Or, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(And, "&&")},
			//			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SemiColon, ";")},
			{regexp.MustCompile(`:`), defaultHandler(Colon, ":")},
			//			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(Comma, ",")},
			{regexp.MustCompile(`\+`), defaultHandler(Plus, "+")},
			{regexp.MustCompile(`-`), defaultHandler(Minus, "-")},
			{regexp.MustCompile(`/`), defaultHandler(Divide, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(Star, "*")},
			//			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

var reservedLu map[string]TokenKind = map[string]TokenKind{
	// "true":    TRUE,
	// "false":   FALSE,
	// "null":    NULL,
	"let":   Let,
	"const": Const,
	// "":   CLASS,
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

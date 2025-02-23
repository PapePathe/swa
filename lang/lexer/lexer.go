package lexer

import (
	"regexp"
)

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

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

func (lex *Lexer) remainder() string {
	return lex.source[lex.position:]
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) atEOF() bool {
	return lex.position >= len(lex.source)
}

type Dialect interface {
	Const() RegexpPattern
	Else() RegexpPattern
	False() RegexpPattern
	For() RegexpPattern
	Function() RegexpPattern
	If() RegexpPattern
	Let() RegexpPattern
	Null() RegexpPattern
	Print() RegexpPattern
	Play() RegexpPattern
	Return() RegexpPattern
	Struct() RegexpPattern
	True() RegexpPattern
	TypeInteger() RegexpPattern
	TypeDecimal() RegexpPattern
	TypeString() RegexpPattern
	Reserved() map[string]TokenKind
	While() RegexpPattern
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
			{regexp.MustCompile(`==`), defaultHandler(Equals, "==")},
			{regexp.MustCompile(`=`), defaultHandler(Assignment, "=")},
			{regexp.MustCompile(`!`), defaultHandler(Not, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LessThanEquals, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LessThan, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GreaterThanEquals, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GreaterThan, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(Or, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(And, "&&")},
			{regexp.MustCompile(`\.`), defaultHandler(Dot, ".")},
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

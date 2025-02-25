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

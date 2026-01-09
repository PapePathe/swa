package lexer

import "regexp"

var dialects = map[string]Dialect{
	"wolof":   Wolof{},
	"malinke": Malinke{},
	"english": English{},
	"french":  French{},
}

type Dialect interface {
	Error(key string, args ...any) error
	DetectionPattern() *regexp.Regexp
	Reserved() map[string]TokenKind
}

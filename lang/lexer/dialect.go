package lexer

import "regexp"

var dialects = map[string]Dialect{
	"wolof":   Wolof{},
	"malinke": Malinke{},
	"english": English{},
	"french":  French{},
	"soussou": Soussou{},
}

type Dialect interface {
	Error(key string, args ...any) error
	// Name returns the name of the dialect as used
	// in the source code (e.g. "english", "français")
	Name() string
	DetectionPattern() *regexp.Regexp
	Reserved() map[string]TokenKind
}

package lexer

var dialects = map[string]Dialect{
	"wolof":   Wolof{},
	"malinke": Malinke{},
	"english": English{},
	"french":  French{},
}

type Dialect interface {
	Patterns() []RegexpPattern
	Reserved() map[string]TokenKind
}

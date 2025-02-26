package lexer

var dialects = map[string]Dialect{
	"wolof":   Wolof{},
	"malinke": Malinke{},
}

type Dialect interface {
	Patterns() []RegexpPattern
}

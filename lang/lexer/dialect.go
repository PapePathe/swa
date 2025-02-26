package lexer

var dialects = map[string]Dialect{
	"wolof":   Wolof{},
	"malinke": Malinke{},
	"english": English{},
}

type Dialect interface {
	Patterns() []RegexpPattern
}

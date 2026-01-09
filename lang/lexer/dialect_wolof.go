package lexer

import (
	"fmt"
	"regexp"

	"swahili/lang/errmsg"
)

type Wolof struct{}

var _ Dialect = (*Wolof)(nil)

func (m Wolof) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialect:wolof;`)
}

func (m Wolof) Patterns() []RegexpPattern {
	return []RegexpPattern{}
}

func (m Wolof) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		return fmt.Errorf("key %s does not exist in wolof dialect translations", key)
	}

	return errmsg.NewAstError(formatted, args)
}

func (m Wolof) Reserved() map[string]TokenKind {
	// TODO: add reserved for wolof
	return map[string]TokenKind{}
}

func (m Wolof) translations() map[string]string {
	// TODO: add translations for error messages
	return map[string]string{}
}

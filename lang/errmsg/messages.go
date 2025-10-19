package errmsg

import "fmt"

func NewAstError(format string, args ...any) error {
	return AstError{Message: format, Args: args}
}

type AstError struct {
	Message string
	Args    []any
}

func (e AstError) Error() string {
	if len(e.Args) > 0 {
		return fmt.Sprintf(e.Message, e.Args[:]...)
	}

	return e.Message
}

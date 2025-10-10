package errmsg

import "fmt"

var formats = map[string]string{
	"fr.ArrayAccessExpression.NameNotASymbol": "L'expression %v n'est pas un nom de variable valide",
	"en.ArrayAccessExpression.NameNotASymbol": "The expression %v is not a correct variable name",
}

func NewAstError(format string, args ...any) AstError {
	return AstError{Message: format, Args: args}
}

type AstError struct {
	Message string
	Args    []any
}

func (e AstError) Error() string {
	if len(e.Args) > 0 {
		return fmt.Sprintf(e.Message, e.Args)
	}

	return e.Message
}

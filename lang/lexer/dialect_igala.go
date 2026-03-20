package lexer

import (
	"regexp"
)

// Country: Nigeria
// Igala is a Volta-Niger language spoken by the Igala people of Kogi State in North-Central
// Nigeria. It is closely related to Yoruba and Itsekiri, forming part of a historical
// linguistic continuum along the Niger River.
type Igala struct{}

var _ Dialect = (*Igala)(nil)

func (m Igala) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`igala:igala;`)
}

func (m Igala) Name() string {
	return "igala"
}

func (m Igala) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ne":           Let,
		"tabat":        Const,
		"alu-ne":       KeywordIf,
		"ne-menye":     KeywordElse,
		"egbe-ne":      KeywordWhile,
		"do":           Function,
		"tro":          Return,
		"start":        Main,
		"ogba":         Print,
		"nhazi":        Struct,
		"ojonogecha":   True,
		"ojonogecha-i": False,
		"bool":         TypeBool,
		"byte":         TypeByte,
		"decimal":      TypeFloat,
		"namba":        TypeInt,
		"namba64":      TypeInt64,
		"eta":          TypeString,
		"obo":          TypeError,
	}
}

func (m Igala) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}

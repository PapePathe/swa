package lexer

type Dialect interface {
	Const() RegexpPattern
	Else() RegexpPattern
	False() RegexpPattern
	For() RegexpPattern
	Function() RegexpPattern
	If() RegexpPattern
	Let() RegexpPattern
	Null() RegexpPattern
	Print() RegexpPattern
	Play() RegexpPattern
	Return() RegexpPattern
	Struct() RegexpPattern
	True() RegexpPattern
	TypeInteger() RegexpPattern
	TypeDecimal() RegexpPattern
	TypeString() RegexpPattern
	Reserved() map[string]TokenKind
	While() RegexpPattern
}

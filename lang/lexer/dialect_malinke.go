package lexer

import "regexp"

type Malinke struct{}

var _ Dialect = (*Malinke)(nil)

func (m Malinke) Null() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`fufafu`),
		handler: defaultHandler(KeywordIf, "fufafu"),
	}
}

func (m Malinke) Print() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`masen`),
		handler: defaultHandler(KeywordIf, "masen"),
	}
}

func (m Malinke) True() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`nondi`),
		handler: defaultHandler(KeywordIf, "nondi"),
	}
}

func (m Malinke) Return() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`ragbilen`),
		handler: defaultHandler(KeywordIf, "ragbilen"),
	}
}

func (m Malinke) Play() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`play_audio`),
		handler: defaultHandler(KeywordIf, "play_audio"),
	}
}

func (m Malinke) False() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`wulé`),
		handler: defaultHandler(KeywordIf, "wulé"),
	}
}

func (m Malinke) If() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`xa`),
		handler: defaultHandler(KeywordIf, "xa"),
	}
}

func (m Malinke) Else() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`xamuara`),
		handler: defaultHandler(KeywordIf, "xamuara"),
	}
}

func (m Malinke) While() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`xaali`),
		handler: defaultHandler(KeywordIf, "xaali"),
	}
}

func (m Malinke) For() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`gbéecé`),
		handler: defaultHandler(KeywordIf, "gbéecé"),
	}
}

func (m Malinke) Const() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`sendendé`),
		handler: defaultHandler(KeywordIf, "sendendé"),
	}
}

func (m Malinke) Let() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`kouicé`),
		handler: defaultHandler(KeywordIf, "kouicé"),
	}
}

func (m Malinke) TypeDecimal() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`decimal`),
		handler: defaultHandler(KeywordIf, "decimal"),
	}
}

func (m Malinke) TypeInteger() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`konti`),
		handler: defaultHandler(KeywordIf, "konti"),
	}
}

func (m Malinke) TypeString() RegexpPattern {
	return RegexpPattern{
		regex:   regexp.MustCompile(`sèbèli`),
		handler: defaultHandler(KeywordIf, "sèbèli"),
	}
}
func (m Malinke) Struct() RegexpPattern {
	return RegexpPattern{}
}

func (m Malinke) Function() RegexpPattern {
	return RegexpPattern{}
}

func (m Malinke) Reserved() map[string]TokenKind {
	return map[string]TokenKind{}
}

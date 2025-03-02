package lexer

import "testing"

func TestTokenize(t *testing.T) {
	result := Tokenize(`
  	dialect:french;
    si(x>0) {
      width = 100;
      height = 100 + 400 - width;
    } sinon {
    	width += 300;
    }
	`)

	expected := []Token{
		{Kind: DialectDeclaration, Value: "dialect"},
		{Kind: Colon, Value: ":"},
		{Kind: Identifier, Value: "french"},
		{Kind: SemiColon, Value: ";"},
		{Kind: KeywordIf, Value: "si"},
		{Kind: OpenParen, Value: "("},
		{Kind: Identifier, Value: "x"},
		{Kind: GreaterThan, Value: ">"},
		{Kind: Number, Value: "0"},
		{Kind: CloseParen, Value: ")"},
		{Kind: OpenCurly, Value: "{"},
		{Kind: Identifier, Value: "width"},
		{Kind: Assignment, Value: "="},
		{Kind: Number, Value: "100"},
		{Kind: SemiColon, Value: ";"},
		{Kind: Identifier, Value: "height"},
		{Kind: Assignment, Value: "="},
		{Kind: Number, Value: "100"},
		{Kind: Plus, Value: "+"},
		{Kind: Number, Value: "400"},
		{Kind: Minus, Value: "-"},
		{Kind: Identifier, Value: "width"},
		{Kind: SemiColon, Value: ";"},
		{Kind: CloseCurly, Value: "}"},
		{Kind: KeywordElse, Value: "sinon"},
		{Kind: OpenCurly, Value: "{"},
		{Kind: Identifier, Value: "width"},
		{Kind: PlusEquals, Value: "+="},
		{Kind: Number, Value: "300"},
		{Kind: SemiColon, Value: ";"},
		{Kind: CloseCurly, Value: "}"},
		{Kind: EOF, Value: "EOF"},
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %s to eq %s at index %d", expected[i], v, i)
		}
	}
}

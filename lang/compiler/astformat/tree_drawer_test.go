package astformat

import (
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"testing"
)

func TestTreeDrawerMultiReturn(t *testing.T) {
	input := `dialect:english;

func divide(a: int, b: int) (int, error) {
    if (b == 0) {
        return 0, error("div by zero"); 
    }
    return a / b, 0; 
}

start() int {
    (res, err) = divide(10, 2);
    print(res);
    return 0;
}
`
	lex, _, _ := lexer.NewFastLexer(input)
	tokens, _ := lex.GetAllTokens()
	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	_ = Draw(&ast)
}

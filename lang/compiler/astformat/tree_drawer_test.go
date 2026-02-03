package astformat 

import (
	"os"
	"path/filepath"
	"strings"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"testing"
)

func TestTreeDrawerAllSwaFiles(t *testing.T) {
	testsDir := "../../../tests"

	err := filepath.Walk(testsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".swa" {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", path, err)
			}

			if len(strings.TrimSpace(string(content))) == 0 {
				t.Skip("Empty file")
			}

			// Tokenize without os.Exit
			lex, _, err := lexer.NewFastLexer(string(content))
			if err != nil {
				t.Skipf("Failed to create lexer (probably missing dialect): %v", err)
			}
			tokens, err := lex.GetAllTokens()
			if err != nil {
				t.Skipf("Failed to tokenize: %v", err)
			}

			// Parse
			ast, err := parser.Parse(tokens)
			if err != nil {
				t.Skipf("Failed to parse %s: %v", path, err)
			}

			// Draw (should not panic)
			_ = Draw(&ast)
		})

		return nil
	})

	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}
}

func TestTreeDrawerMultiReturn(t *testing.T) {
	input := `dialect:english;

func divide(a: int, b: int) (int, error) {
    if (b == 0) {
        return 0, error; 
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

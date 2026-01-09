package lexer

import (
	"fmt"
	"testing"
)

// TestFastasic tests basic token recognition
func TestFastasic(t *testing.T) {
	source := `dialect:english;
let x = 42;
print "Hello, World!";`

	fastLex, _, err := NewFastLexer(source)
	if err != nil {
		t.Fatalf("Failed to create  %v", err)
	}

	tokens, err := fastLex.GetAllTokens()
	if err != nil {
		t.Fatalf("Lexing failed: %v", err)
	}

	expectedTokens := []struct {
		kind   TokenKind
		lexeme string
	}{
		{DialectDeclaration, "dialect"},
		{Colon, ":"},
		{Identifier, "english"},
		{SemiColon, ";"},
		{Let, "let"},
		{Identifier, "x"},
		{Assignment, "="},
		{Number, "42"},
		{SemiColon, ";"},
		{Print, "print"},
		{String, "Hello, World!"},
		{SemiColon, ";"},
	}

	if len(tokens) != len(expectedTokens) {
		t.Errorf("Expected %d tokens, got %d", len(expectedTokens), len(tokens))
		t.Logf("Tokens %v", tokens)
		t.Logf("ExpectedTokens %v", expectedTokens)
	}

	for i, expected := range expectedTokens {
		if i >= len(tokens) {
			break
		}
		if tokens[i].Kind != expected.kind {
			t.Errorf("Token %d: expected kind %v, got %v", i, expected.kind, tokens[i].Kind)
		}

		if tokens[i].Value != expected.lexeme {
			t.Errorf("Token %d: expected lexeme %q, got %q", i, expected.lexeme, tokens[i].Value)
		}
	}
}

// TestFastperators tests all operators
func TestFastperators(t *testing.T) {
	source := `dialect:english;
x = y + z * 2 / 3 % 4;
a == b != c < d > e <= f >= g;
!condition && (a || b);
x += 1;`

	fastLex, _, err := NewFastLexer(source)
	if err != nil {
		t.Fatalf("Failed to create  %v", err)
	}

	tokens, err := fastLex.GetAllTokens()
	if err != nil {
		t.Fatalf("Lexing failed: %v", err)
	}

	// Verify key operators are present
	operatorTests := []struct {
		lexeme string
		kind   TokenKind
	}{
		{"=", Assignment},
		{"+", Plus},
		{"*", Star},
		{"/", Divide},
		{"%", Modulo},
		{"==", Equals},
		{"!=", NotEquals},
		{"<", LessThan},
		{">", GreaterThan},
		{"<=", LessThanEquals},
		{">=", GreaterThanEquals},
		{"!", Not},
		{"&&", And},
		{"||", Or},
		{"+=", PlusEquals},
	}

	for _, test := range operatorTests {
		found := false
		for _, token := range tokens {
			if token.Value == test.lexeme && token.Kind == test.kind {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Operator %s (kind %v) not found in tokens", test.lexeme, test.kind)
		}
	}
}

// TestFastumbers tests number lexing
func TestFastumbers(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []string
		kind     TokenKind
	}{
		{
			name:     "integers",
			source:   "0 42 -123 999",
			expected: []string{"0", "42", "-123", "999"},
			kind:     Number,
		},
		{
			name:     "floats",
			source:   "3.14 -2.718 0.0 100.001",
			expected: []string{"3.14", "-2.718", "0.0", "100.001"},
			kind:     Float,
		},
		{
			name:     "mixed",
			source:   "x = 42 + 3.14",
			expected: []string{"42", "3.14"},
			kind:     Number, // Will check both kinds
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := fmt.Sprintf("dialect:english;\n%s", tt.source)
			fastLex, _, err := NewFastLexer(src)
			if err != nil {
				t.Fatalf("Failed to create  %v", err)
			}

			tokens, err := fastLex.GetAllTokens()
			if err != nil {
				t.Fatalf("Lexing failed: %v", err)
			}

			// Collect number tokens
			var numberTokens []Token
			for _, token := range tokens {
				if token.Kind == Number || token.Kind == Float {
					numberTokens = append(numberTokens, token)
				}
			}

			if len(numberTokens) != len(tt.expected) {
				t.Errorf("Expected %d number tokens, got %d", len(tt.expected), len(numberTokens))
			}

			for i, expected := range tt.expected {
				if i >= len(numberTokens) {
					break
				}
				if numberTokens[i].Value != expected {
					t.Errorf("Token %d: expected %q, got %q", i, expected, numberTokens[i].Value)
				}
			}
		})
	}
}

// TestFasttrings tests string lexing with escape sequences
func TestFasttrings(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expected    string
		shouldError bool
	}{
		{
			name:     "simple string",
			source:   `"Hello, World!"`,
			expected: "Hello, World!",
		},
		{
			name:     "string with escape",
			source:   `"Line1\nLine2\tTab"`,
			expected: "Line1\nLine2\tTab",
		},
		{
			name:     "empty string",
			source:   `""`,
			expected: "",
		},
		{
			name:        "unterminated string",
			source:      `"Hello, World!`,
			shouldError: true,
		},
		{
			name:     "string with quotes",
			source:   `"He said, \"Hello!\""`,
			expected: `He said, "Hello!"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := fmt.Sprintf("dialect:english;\n%s", tt.source)
			fastLex, _, err := NewFastLexer(src)
			if err != nil {
				t.Fatalf("Failed to create  %v", err)
			}

			tokens, err := fastLex.GetAllTokens()

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Lexing failed: %v", err)
			}

			if len(tokens) != 5 {
				t.Errorf("Expected 5 tokens, got %d", len(tokens))
				return
			}

			if tokens[4].Kind != String {
				t.Errorf("Expected String, got %v", tokens[0].Kind)
			}

			if tokens[4].Value != tt.expected {
				t.Errorf("Expected lexeme %q, got %q", tt.expected, tokens[0].Value)
			}
		})
	}
}

// TestFastharacters tests character lexing
func TestFastharacters(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expected    string
		shouldError bool
	}{
		{
			name:     "simple char",
			source:   `'a'`,
			expected: "a",
		},
		{
			name:     "escape char",
			source:   `'\n'`,
			expected: "\n",
		},
		{
			name:     "tab char",
			source:   `'\t'`,
			expected: "\t",
		},
		{
			name:     "quote char",
			source:   `'\''`,
			expected: "'",
		},
		{
			name:        "unterminated char",
			source:      `'a`,
			shouldError: true,
		},
		{
			name:        "empty char",
			source:      `''`,
			shouldError: true,
		},
		{
			name:        "invalid escape",
			source:      `'\x'`,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := fmt.Sprintf("dialect:english;\n%s", tt.source)
			fastLex, _, err := NewFastLexer(src)
			if err != nil {
				t.Fatalf("Failed to create  %v", err)
			}

			tokens, err := fastLex.GetAllTokens()

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Lexing failed: %v", err)
			}

			if len(tokens) != 5 {
				t.Errorf("Expected 1 token, got %d", len(tokens))
				return
			}

			if tokens[4].Kind != Character {
				t.Errorf("Expected Character, got %v", tokens[0].Kind)
			}

			if tokens[4].Value != tt.expected {
				t.Errorf("Expected lexeme %q, got %q", tt.expected, tokens[0].Value)
			}
		})
	}
}

// TestFastdentifiers tests identifier lexing
func TestFastdentifiers(t *testing.T) {
	tests := []struct {
		name   string
		source string
		valid  bool
	}{
		{"simple", "x", true},
		{"underscore", "_private", true},
		{"camelCase", "myVariable", true},
		{"with numbers", "var123", true},
		{"starting number", "123var", false}, // Will be lexed as number + identifier
		{"special chars", "my-var", false},   // Minus will be separate operator
		{"unicode", "café", true},
		{"multiple", "x y z123 _test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := fmt.Sprintf("dialect:english;\n%s", tt.source)
			fastLex, _, err := NewFastLexer(src)
			if err != nil {
				t.Fatalf("Failed to create  %v", err)
			}

			tokens, err := fastLex.GetAllTokens()
			if err != nil {
				t.Fatalf("Lexing failed: %v", err)
			}

			// For valid identifiers, check we got identifier tokens
			if tt.valid {
				hasIdentifier := false
				for _, token := range tokens {
					if token.Kind == Identifier {
						hasIdentifier = true
						break
					}
				}
				if !hasIdentifier {
					t.Errorf("Expected at least one identifier token")
				}
			}
		})
	}
}

// TestFasteywords tests keyword recognition
func TestFasteywords(t *testing.T) {
	source := `dialect:english; if else while func struct let const int float string return print`

	fastLex, _, err := NewFastLexer(source)
	if err != nil {
		t.Fatalf("Failed to create  %v", err)
	}

	tokens, err := fastLex.GetAllTokens()
	if err != nil {
		t.Fatalf("Lexing failed: %v", err)
	}

	expectedKeywords := []TokenKind{
		DialectDeclaration,
		Colon,
		Identifier,
		SemiColon,
		KeywordIf,
		KeywordElse,
		KeywordWhile,
		Function,
		Struct,
		Let,
		Const,
		TypeInt,
		TypeFloat,
		TypeString,
		Return,
		Print,
	}

	if len(tokens) != len(expectedKeywords) {
		t.Errorf("Expected %d tokens, got %d", len(expectedKeywords), len(tokens))
	}

	for i, token := range tokens {
		if i >= len(expectedKeywords) {
			break
		}
		if token.Kind != expectedKeywords[i] {
			t.Errorf("Token %d: expected %v, got %v", i, expectedKeywords[i], token.Kind)
		}
	}
}

// TestFastrench tests French dialect
func TestFastrench(t *testing.T) {
	source := `dialecte:français;
variable x = 42;
afficher "Bonjour!";
tantque vrai {
    afficher "En boucle";
}`

	fastLex, _, err := NewFastLexer(source)
	if err != nil {
		t.Fatalf("Failed to create  %v", err)
	}

	//	if dialect.DetectionKeyword() != "dialecte:français;" {
	//		t.Errorf("Expected French dialect, got %v", dialect)
	//	}

	tokens, err := fastLex.GetAllTokens()
	if err != nil {
		t.Fatalf("Lexing failed: %v", err)
	}

	// Check for French keywords
	expectedKeywords := []struct {
		kind   TokenKind
		lexeme string
	}{
		{DialectDeclaration, "dialecte"},
		{Colon, ":"},
		{Identifier, "français"},
		{SemiColon, ";"},
		{Let, "variable"},
		{Identifier, "x"},
		{Assignment, "="},
		{Number, "42"},
		{SemiColon, ";"},
		{Print, "afficher"},
		{String, "Bonjour!"},
		{SemiColon, ";"},
		{KeywordWhile, "tantque"},
		{Identifier, "vrai"},
		{OpenCurly, "{"},
		{Print, "afficher"},
		{String, "En boucle"},
		{SemiColon, ";"},
		{CloseCurly, "}"},
	}

	for i, expected := range expectedKeywords {
		if i >= len(tokens) {
			t.Errorf("Missing token %d: expected %v", i, expected)
			continue
		}
		if tokens[i].Kind != expected.kind {
			t.Errorf("Token %d: expected kind %v, got %v", i, expected.kind, tokens[i].Kind)
		}
		if tokens[i].Value != expected.lexeme {
			t.Errorf("Token %d: expected lexeme %q, got %q", i, expected.lexeme, tokens[i].Value)
		}
	}
}

// TestFastomments tests comment handling
func TestFastomments(t *testing.T) {
	source := `// This is a comment
x = 42; // Assign value
// Another comment`

	fastLex, err := NewFastLexerWithDialect(source, English{})
	if err != nil {
		t.Fatalf("Failed to create  %v", err)
	}

	tokens, err := fastLex.GetAllTokens()
	if err != nil {
		t.Fatalf("Lexing failed: %v", err)
	}

	// Comments should be skipped, so we should only get:
	// x = 42;
	expectedTokens := 4 // x, =, 42, ;

	if len(tokens) != expectedTokens {
		t.Errorf("Expected %d tokens (comments skipped), got %d", expectedTokens, len(tokens))
		for i, tok := range tokens {
			t.Logf("Token %d: %v", i, tok)
		}
	}
}

// TestFasthitespace tests whitespace handling
func TestFasthitespace(t *testing.T) {
	source := "x\t=\t42\n\ny  =   \t100  "

	fastLex, err := NewFastLexerWithDialect(source, English{})
	if err != nil {
		t.Fatalf("Failed to create  %v", err)
	}

	tokens, err := fastLex.GetAllTokens()
	if err != nil {
		t.Fatalf("Lexing failed: %v", err)
	}

	// Should have: x, =, 42, y, =, 100
	expectedTokens := 6
	if len(tokens) != expectedTokens {
		t.Errorf("Expected %d tokens, got %d", expectedTokens, len(tokens))
		for i, tok := range tokens {
			t.Logf("Token %d: %v", i, tok)
		}
	}
}

// TestFastdgeCases tests edge cases
func TestFastdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{"empty", ""},
		{"only whitespace", "   \t\n  "},
		{"only comment", "// Just a comment"},
		{"unicode", "café = '☕'"},
		{"nested braces", "{{{}}}[[[]]]((()))"},
		{"complex expression", "result = (a + b) * (c - d) / e % f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := fmt.Sprintf("dialect:english;\n%s", tt.source)
			fastLex, _, err := NewFastLexer(src)
			if err != nil {
				t.Fatalf("Failed to create  %v", err)
			}

			_, err = fastLex.GetAllTokens()
			if err != nil {
				t.Errorf("Lexing failed for %q: %v", tt.name, err)
			}
		})
	}
}

// TestFastineNumbers tests line and column tracking
func TestFastineNumbers(t *testing.T) {
	source := `line1
  line2
    line3`

	fastLex, err := NewFastLexerWithDialect(source, English{})
	if err != nil {
		t.Fatalf("Failed to create  %v", err)
	}

	tokens, err := fastLex.GetAllTokens()
	if err != nil {
		t.Fatalf("Lexing failed: %v", err)
	}

	// Check line numbers
	expectedLines := []int{1, 2, 3, 3, 4, 4}
	expectedColumns := []int{1, 3, 5, 8, 5, 10}

	for i, token := range tokens {
		if i >= len(expectedLines) {
			break
		}
		if token.Line != expectedLines[i] {
			t.Errorf("Token %d (%s): expected line %d, got %d",
				i, token.Value, expectedLines[i], token.Line)
		}
		if token.Column != expectedColumns[i] {
			t.Errorf("Token %d (%s): expected column %d, got %d",
				i, token.Value, expectedColumns[i], token.Column)
		}
	}
}

// BenchmarkFastbenchmarks the performance
func BenchmarkFast(b *testing.B) {
	source := `dialect:english;
func factorial(n int) int {
    if n <= 1 {
        return 1;
    }
    return n * factorial(n - 1);
}

start {
    result = factorial(5);
    print "Factorial of 5 is: " + result;
}`

	for i := 0; i < b.N; i++ {
		fastLex, _, err := NewFastLexer(source)
		if err != nil {
			b.Fatalf("Failed to create  %v", err)
		}

		_, err = fastLex.GetAllTokens()
		if err != nil {
			b.Fatalf("Lexing failed: %v", err)
		}
	}
}

// BenchmarkFastrench benchmarks French dialect lexing
func BenchmarkFastrench(b *testing.B) {
	source := `dialecte:français;
fonction factoriel(n entier) entier {
    si n <= 1 {
        retourner 1;
    }
    retourner n * factoriel(n - 1);
}

demarrer {
    résultat = factoriel(5);
    afficher "Le factoriel de 5 est: " + résultat;
}`

	for i := 0; i < b.N; i++ {
		fastLex, _, err := NewFastLexer(source)
		if err != nil {
			b.Fatalf("Failed to create  %v", err)
		}

		_, err = fastLex.GetAllTokens()
		if err != nil {
			b.Fatalf("Lexing failed: %v", err)
		}
	}
}

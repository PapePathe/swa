package lexer_test

import (
	"swahili/lang/lexer"
	"testing"
)

func TestNewToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		kind  lexer.TokenKind
		value string
		want  lexer.Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lexer.NewToken(tt.kind, tt.value)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

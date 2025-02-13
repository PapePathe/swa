package main

import (
	"os"
	"swahili/lang/lexer"
)

func main() {
	bytes, _ := os.ReadFile("./examples/malinke.swa")
	source := string(bytes)

	tokens := lexer.Tokenize(source)

	for _, token := range tokens {
		token.Debug()
	}
}

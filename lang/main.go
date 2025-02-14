package main

import (
	"os"
	"swahili/lang/lexer"
	"swahili/lang/parser"

	"github.com/sanity-io/litter"
)

func main() {
	bytes, _ := os.ReadFile("./examples/malinke/age_calculator.swa")
	source := string(bytes)

	tokens := lexer.Tokenize(source)
	st := parser.Parse(tokens)
	litter.Dump(st)

}

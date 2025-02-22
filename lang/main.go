package main

import (
	"os"

	"github.com/sanity-io/litter"

	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func main() {
	file := "./examples/malinke/structs.swa"
	bytes, _ := os.ReadFile(file)
	source := string(bytes)
	tokens := lexer.Tokenize(source)
	st := parser.Parse(tokens)

	litter.Dump(st)
}

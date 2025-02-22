package main

import (
	"fmt"
	"os"

	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func main() {
	file := "./examples/malinke/structs.swa"
	bytes, _ := os.ReadFile(file)
	source := string(bytes)
	tokens := lexer.Tokenize(source)
	st := parser.Parse(tokens)

	fmt.Println(st)
}

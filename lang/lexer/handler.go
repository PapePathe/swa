package lexer

import (
	"regexp"
	"strconv"
)

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *Lexer, _ *regexp.Regexp) {
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value, lex.line))
	}
}

func commentHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	if match != nil {
		lex.advanceN(match[1])
	}
}

func characterHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	charLiteral := lex.remainder()[match[0]:match[1]]

	lex.push(NewToken(Character, charLiteral, lex.line))
	lex.advanceN(match[1])
}

func stringHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]

	text, err := strconv.Unquote(stringLiteral)
	if err != nil {
		panic(err)
	}

	newToken := NewToken(String, text, lex.line)

	lex.push(newToken)
	lex.advanceN(len(stringLiteral))
}

func newlineHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.advanceN(len(match))
	lex.newLine()
}

func floatHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(Float, match, lex.line))
	lex.advanceN(len(match))
}

func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(Number, match, lex.line))
	lex.advanceN(len(match))
}

func symbolHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if kind, found := lex.reservedWords[match]; found {
		lex.push(NewToken(kind, match, lex.line))
	} else {
		lex.push(NewToken(Identifier, match, lex.line))
	}

	lex.advanceN(len(match))
}

func skipHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())

	lex.advanceN(match[1])
}

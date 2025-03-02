package lexer

import (
	"regexp"
)

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *Lexer, _ *regexp.Regexp) {
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value))
	}
}

func commentHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	if match != nil {
		lex.advanceN(match[1])
	}
}

func stringHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]

	lex.push(NewToken(String, stringLiteral))
	lex.advanceN(len(stringLiteral))
}

func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(Number, match))
	lex.advanceN(len(match))
}

func symbolHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	if kind, found := lex.reservedWords[match]; found {
		lex.push(NewToken(kind, match))
	} else {
		lex.push(NewToken(Identifier, match))
	}

	lex.advanceN(len(match))
}

func skipHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

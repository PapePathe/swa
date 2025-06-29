/*
* swahili/lang
* Copyright (C) 2025  Papa Pathe SENE
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package lexer

// TokenKind is a custom type representing the kind of a token.
type TokenKind int

var tks = map[TokenKind]string{
	And:                "AND",
	Assignment:         "ASSIGNMENT",
	Colon:              "COLON",
	Comma:              "COMMA",
	CloseCurly:         "CLOSE_CURLY",
	CloseParen:         "CLOSE_PAREN",
	CloseBracket:       "CLOSE_BRACKET",
	Const:              "CONST",
	DialectDeclaration: "DIALECT",
	Divide:             "DIVIDE",
	Dot:                "DOT",
	EOF:                "EOF",
	Equals:             "EQUALS",
	Function:           "FUNCTION",
	GreaterThan:        "GREATER_THAN",
	GreaterThanEquals:  "GREATER_THAN_EQUALS",
	Identifier:         "IDENTIFIER",
	KeywordIf:          "IF",
	KeywordElse:        "ELSE",
	LessThan:           "LESS_THAN",
	LessThanEquals:     "LESS_THAN_EQUALS",
	Let:                "LET",
	Main:               "MAIN_PROGRAM",
	Minus:              "MINUS",
	Multiply:           "MULTIPLY",
	Not:                "NOT",
	NotEquals:          "NOT_EQUALS",
	Number:             "NUMBER",
	OpenCurly:          "OPEN_CURLY",
	OpenParen:          "OPEN_PAREN",
	OpenBracket:        "OPEN_BRACKET",
	Or:                 "OR",
	Plus:               "PLUS",
	PlusEquals:         "PLUS_EQUAL",
	Print:              "PRINT",
	QuestionMark:       "QUESTION_MARK",
	SemiColon:          "SEMI_COLON",
	Return:             "RETURN",
	StarEquals:         "STAR_EQUALS",
	Struct:             "STRUCT",
	String:             "STRING",
	Star:               "STAR",
	TypeInt:            "INT",
}

// String s a string representation of the TokenKind.
func (k TokenKind) String() string {
	switch k {
	case Character:
		return "CHARACTER"
	case Print:
		return "PRINT"
	case DialectDeclaration:
		return "DIALECT"
	case StarEquals:
		return "STAR_EQUALS"
	case Dot:
		return "DOT"
	case Struct:
		return "STRUCT"
	case Let:
		return "LET"
	case Const:
		return "CONST"
	case Or:
		return "OR"
	case And:
		return "AND"
	case KeywordIf:
		return "IF"
	case KeywordElse:
		return "ELSE"
	case TypeInt:
		return "INT"
	case Identifier:
		return "IDENTIFIER"
	case Assignment:
		return "ASSIGNMENT"
	case String:
		return "STRING"
	case Colon:
		return "COLON"
	case Comma:
		return "COMMA"
	case CloseCurly:
		return "CLOSE_CURLY"
	case CloseParen:
		return "CLOSE_PAREN"
	case CloseBracket:
		return "CLOSE_BRACKET"
	case Divide:
		return "DIVIDE"
	case EOF:
		return "EOF"
	case Equals:
		return "EQUALS"
	case GreaterThan:
		return "GREATER_THAN"
	case GreaterThanEquals:
		return "GREATER_THAN_EQUALS"
	case LessThan:
		return "LESS_THAN"
	case LessThanEquals:
		return "LESS_THAN_EQUALS"
	case Minus:
		return "MINUS"
	case Multiply:
		return "MULTIPLY"
	case Not:
		return "NOT"
	case NotEquals:
		return "NOT_EQUALS"
	case OpenCurly:
		return "OPEN_CURLY"
	case OpenParen:
		return "OPEN_PAREN"
	case OpenBracket:
		return "OPEN_BRACKET"
	case Plus:
		return "PLUS"
	case PlusEquals:
		return "PLUS_EQUAL"
	case Star:
		return "STAR"
	case SemiColon:
		return "SEMI_COLON"
	case QuestionMark:
		return "QUESTION_MARK"
	case Number:
		return "NUMBER"
	}

	return ""
}

const (
	// Let represents the "let" keyword.
	Let TokenKind = iota
	// Const represents the "const" keyword.
	Const
	// Or represents the "or" logical operator.
	Or
	// And represents the "and" logical operator.
	And
	// KeywordIf represents the "if" keyword.
	KeywordIf
	// KeywordElse represents the "else" keyword.
	KeywordElse

	// TypeInt represents the "int" type.
	TypeInt
	// String represents a string literal.
	String
	// Number represents a numeric literal.
	Number

	// Assignment represents an assignment operator.
	Assignment
	// Colon represents a colon symbol.
	Colon
	// Comma represents a comma symbol.
	Comma
	// CloseCurly represents a closing curly brace symbol.
	CloseCurly
	// CloseParen represents a closing parenthesis symbol.
	CloseParen
	// CloseBracket represents a closing bracket symbol.
	CloseBracket
	// Divide represents the divide operator.
	Divide
	// EOF represents the end of file token.
	EOF
	// Equals represents the equality operator.
	Equals
	// GreaterThan represents the greater than operator.
	GreaterThan
	// GreaterThanEquals represents the greater than or equal to operator.
	GreaterThanEquals
	// LessThan represents the less than operator.
	LessThan
	// LessThanEquals represents the less than or equal to operator.
	LessThanEquals
	// Minus represents the minus operator.
	Minus
	// MinusEquals.
	MinusEquals
	// Multiply represents the multiplication operator.
	Multiply
	// Not represents the negation operator.
	Not
	// NotEquals represents the not equal operator.
	NotEquals
	// OpenCurly represents an opening curly brace symbol.
	OpenCurly
	// OpenParen represents an opening parenthesis symbol.
	OpenParen
	// OpenBracket represents an opening bracket symbol.
	OpenBracket
	// Plus represents the addition operator.
	Plus
	// Plus represents the plus equals operator.
	PlusEquals
	// Star represents the multiplication symbol (star).
	Star
	// SemiColon represents a semicolon symbol.
	SemiColon
	// QuestionMark represents a question mark symbol.
	QuestionMark
	// Identifier represents an identifier (variable or function name).
	Identifier
	// Struct ...
	Struct
	// Dot ...
	Dot
	// StarEquals ...
	StarEquals
	// Dialect ...
	DialectDeclaration
	// Print
	Print
	// A single ascii or utf8 character
	Character
	Function
	Main
	Return
)

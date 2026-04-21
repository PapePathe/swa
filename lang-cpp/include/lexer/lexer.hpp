#include <iostream>
#include <vector>
#include <string>
#include <cctype>
#include <assert.h>
#include <iomanip>
#include <unordered_map>

enum class TokenType {
  LET, CONST, OR, AND, KEYWORD_IF, KEYWORD_ELSE, KEYWORD_WHILE,
  TYPE_INT, TYPE_INT64, TYPE_FLOAT, TYPE_STRING, TYPE_BOOL, TYPE_BYTE,
  STRUCT, PRINT, FUNCTION, MAIN, RETURN, TRUE, FALSE,
  IDENTIFIER, STRING, NUMBER, FLOAT, CHARACTER,
  ASSIGNMENT, COLON, COMMA, DOT, SEMICOLON, QUESTION_MARK, TILDE,
  OPEN_CURLY, CLOSE_CURLY, OPEN_PAREN, CLOSE_PAREN, OPEN_BRACKET, CLOSE_BRACKET,
  EQUALS, NOT_EQUALS, GREATER_THAN, GREATER_THAN_EQUALS,
  LESS_THAN, LESS_THAN_EQUALS, PLUS, PLUS_EQUALS, MINUS, MINUS_EQUALS,
  MULTIPLY, STAR_EQUALS, DIVIDE, MODULO, DOUBLE_STAR, AMPERSAND, VARIADIC,
  DIALECT_DECLARATION, ZERO, WHITESPACE, COMMENT, NEWLINE, 
  TYPE_ERROR, END_OF_FILE, UNKNOWN
};

struct Token {
    TokenType type;
    std::string value;
};

class Lexer {
    std::string src;
    std::unordered_map<std::string, TokenType> keywords; 
    size_t pos = 0;

    char peek();
    char get();

  public:
    Lexer(
        std::string s, 
        std::unordered_map<std::string, TokenType> k
    ) : src(s), pos(0), keywords(k)  {}
    std::vector<Token> tokenize();
};

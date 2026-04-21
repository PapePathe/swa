#include <iostream>
#include <vector>
#include <string>
#include <cctype>
#include <assert.h>
#include <iomanip>
#include <unordered_map>
#include "tokentype.hpp"

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

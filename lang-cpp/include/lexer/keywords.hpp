#include <unordered_map>
#include "tokentype.hpp"

const std::unordered_map<std::string, TokenType> KEYWORDS_ENGLISH = {
    {"const",   TokenType::CONST},
    {"dialect", TokenType::DIALECT_DECLARATION}, 
    {"func",    TokenType::FUNCTION}, 
    {"else",    TokenType::KEYWORD_ELSE},
    {"if",      TokenType::KEYWORD_IF}, 
    {"while",   TokenType::KEYWORD_WHILE}, 
    {"let",     TokenType::LET}, 
    {"int",     TokenType::TYPE_INT},
    {"byte",    TokenType::TYPE_BYTE}, 
    {"error",   TokenType::TYPE_ERROR}, 
    {"float",   TokenType::TYPE_FLOAT}, 
    {"string",  TokenType::TYPE_STRING},
    {"bool",    TokenType::TYPE_BOOL}, 
    {"func",    TokenType::FUNCTION},
    {"return",  TokenType::RETURN}, 
    {"and",     TokenType::AND},
    {"or",      TokenType::OR}, 
    {"true",    TokenType::TRUE},
    {"false",   TokenType::FALSE},
    {"zero",    TokenType::ZERO}
};

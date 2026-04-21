#include "lexer/lexer.hpp"
#include <vector>
#include <cctype>
#include <unordered_map>

// Mapping our Go-style keyword list to TokenTypes
const std::unordered_map<std::string, TokenType> KEYWORDS = {
    {"let", TokenType::LET}, {"const", TokenType::CONST},
    {"if", TokenType::KEYWORD_IF}, {"else", TokenType::KEYWORD_ELSE},
    {"while", TokenType::KEYWORD_WHILE}, {"int", TokenType::TYPE_INT},
    {"float", TokenType::TYPE_FLOAT}, {"string", TokenType::TYPE_STRING},
    {"bool", TokenType::TYPE_BOOL}, {"func", TokenType::FUNCTION},
    {"return", TokenType::RETURN}, {"and", TokenType::AND},
    {"or", TokenType::OR}, {"true", TokenType::TRUE},
    {"false", TokenType::FALSE}
};

char Lexer::peek() {
    return pos < src.length() ? src[pos] : '\0';
}

char Lexer::get() {
    return pos < src.length() ? src[pos++] : '\0';
}

std::vector<Token> Lexer::tokenize() {
    std::vector<Token> tokens;

    while (pos < src.length()) {
        char c = peek();

        // 1. Skip Whitespace
        if (isspace(static_cast<unsigned char>(c))) {
            get();
            continue;
        }

        // 2. Numbers (Decimals supported)
        if (isdigit(static_cast<unsigned char>(c))) {
            std::string num;
            bool hasDot = false;
            while (isdigit(peek()) || (peek() == '.' && !hasDot)) {
                if (peek() == '.') hasDot = true;
                num += get();
            }
            tokens.push_back({hasDot ? TokenType::FLOAT : TokenType::NUMBER, num});
        }

        // 3. Identifiers & Keywords (UTF-8 / Accented support)
        else if (isalpha(static_cast<unsigned char>(c)) || c == '_' || static_cast<unsigned char>(c) > 127) {
            std::string ident;
            while (isalnum(static_cast<unsigned char>(peek())) || peek() == '_' || static_cast<unsigned char>(peek()) > 127) {
                ident += get();
            }
            
            if (KEYWORDS.count(ident)) {
                tokens.push_back({KEYWORDS.at(ident), ident});
            } else {
                tokens.push_back({TokenType::IDENTIFIER, ident});
            }
        }

        // 4. Operators & Symbols
        else {
            char current = get();
            switch (current) {
                case '+':
                    if (peek() == '=') { get(); tokens.push_back({TokenType::PLUS_EQUALS, "+="}); }
                    else { tokens.push_back({TokenType::PLUS, "+"}); }
                    break;
                case '-':
                    if (peek() == '=') { get(); tokens.push_back({TokenType::MINUS_EQUALS, "-="}); }
                    else { tokens.push_back({TokenType::MINUS, "-"}); }
                    break;
                case '*':
                    if (peek() == '*') { get(); tokens.push_back({TokenType::DOUBLE_STAR, "**"}); }
                    else if (peek() == '=') { get(); tokens.push_back({TokenType::STAR_EQUALS, "*="}); }
                    else { tokens.push_back({TokenType::MULTIPLY, "*"}); }
                    break;
                case '=':
                    if (peek() == '=') { get(); tokens.push_back({TokenType::EQUALS, "=="}); }
                    else { tokens.push_back({TokenType::ASSIGNMENT, "="}); }
                    break;
                case '&':
                    if (peek() == '&') { get(); tokens.push_back({TokenType::AND, "&&"}); }
                    else { tokens.push_back({TokenType::AMPERSAND, "&"}); }
                    break;
                case '|':
                    if (peek() == '|') { get(); tokens.push_back({TokenType::OR, "||"}); }
                    break;
                case '.':
                    if (peek() == '.') {
                        get(); // consume second .
                        if (peek() == '.') { get(); tokens.push_back({TokenType::VARIADIC, "..."}); }
                    } else {
                        tokens.push_back({TokenType::DOT, "."});
                    }
                    break;
                case '(': tokens.push_back({TokenType::OPEN_PAREN, "("}); break;
                case ')': tokens.push_back({TokenType::CLOSE_PAREN, ")"}); break;
                case '{': tokens.push_back({TokenType::OPEN_CURLY, "{"}); break;
                case '}': tokens.push_back({TokenType::CLOSE_CURLY, "}"}); break;
                case ';': tokens.push_back({TokenType::SEMICOLON, ";"}); break;
                case ',': tokens.push_back({TokenType::COMMA, ","}); break;
                default:
                    tokens.push_back({TokenType::TYPE_ERROR, std::string(1, current)});
                    break;
            }
        }
    }

    tokens.push_back({TokenType::END_OF_FILE, ""});
    return tokens;
};

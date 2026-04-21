#include <iostream>
#include <vector>
#include <string>
#include <cctype>
#include <assert.h>
#include <iomanip>
#include <unordered_map>
#include "lexer/lexer.hpp"
#include "lexer/keywords.hpp"

std::string tokenTypeToString(TokenType type) {
    switch (type) {
        case TokenType::FUNCTION:    return "FUNCTION";
        case TokenType::OR:          return "OR";
        case TokenType::LET:         return "LET";
        case TokenType::NUMBER:      return "NUMBER";
        case TokenType::IDENTIFIER:  return "IDENTIFIER";
        case TokenType::PLUS:        return "PLUS";
        case TokenType::MINUS:       return "MINUS";
        case TokenType::MULTIPLY:    return "MULTIPLY";
        case TokenType::DIVIDE:      return "DIVIDE";
        case TokenType::END_OF_FILE: return "EOF";
        default:                     return "UNKNOWN";
    }
}

int main() {
  while(true) {
        // 1. Get input from the user
    std::string input;
    std::cout << "Enter expression (e.g., total_é + 5 * x): ";
    std::getline(std::cin, input);

    // 2. Initialize the Lexer
    Lexer lexer(input, KEYWORDS_ENGLISH);

    // 3. Generate the tokens
    std::vector<Token> tokens = lexer.tokenize();

    // 4. Process/Display the results
    std::cout << "\n--- Token Stream ---\n";
    std::cout << std::left << std::setw(15) << "TYPE" << "VALUE" << "\n";
    std::cout << std::string(25, '-') << "\n";

    for (const auto& token : tokens) {
        std::cout << std::left << std::setw(15) << tokenTypeToString(token.type) 
                  << "[" << token.value << "]" << "\n";
    }
  }


    return 0;
}

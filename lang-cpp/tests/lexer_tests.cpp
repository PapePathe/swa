#include <gtest/gtest.h>
#include <string>
#include <vector>
#include "lexer/lexer.hpp"
#include "lexer/keywords.hpp"

class LexerTest : public ::testing::Test {
protected:
    std::vector<Token> getTokens(
        const std::string& input, 
        const std::unordered_map<std::string, TokenType> keywords) {
        Lexer lexer(input, keywords);
        return lexer.tokenize();
    }
};

std::string tokenTypeToString(TokenType type) {
    switch (type) {
        case TokenType::DIALECT:       return "DIALECT";
        case TokenType::FUNCTION:      return "FUNCTION";
        case TokenType::CONST:         return "CONST";
        case TokenType::OR:            return "OR";
        case TokenType::LET:           return "LET";
        case TokenType::NUMBER:        return "NUMBER";
        case TokenType::IDENTIFIER:    return "IDENTIFIER";
        case TokenType::PLUS:          return "PLUS";
        case TokenType::MINUS:         return "MINUS";
        case TokenType::MULTIPLY:      return "MULTIPLY";
        case TokenType::DIVIDE:        return "DIVIDE";
        case TokenType::COLON:         return "COLON";
        case TokenType::SEMICOLON:     return "SEMICOLON";
        case TokenType::END_OF_FILE: return "EOF";
        default:                     return "UNKNOWN";
    }
}

void AssertToken(const Token& t, TokenType expectedType, const std::string& expectedValue) {
    EXPECT_EQ(t.type, expectedType) << tokenTypeToString(expectedType) << " do not match " << tokenTypeToString(t.type);
    EXPECT_EQ(t.value, expectedValue);
}

TEST_F(LexerTest, HandlesBasicArithmetic) {
  auto tokens = getTokens("123 + 456", KEYWORDS_ENGLISH);

  ASSERT_EQ(tokens.size(), 4); 
  AssertToken(tokens[0], TokenType::NUMBER,      "123");
  AssertToken(tokens[1], TokenType::PLUS,        "+");
  AssertToken(tokens[2], TokenType::NUMBER,      "456");
  AssertToken(tokens[3], TokenType::END_OF_FILE, "");
} 

TEST_F(LexerTest, HandlesAccentedIdentifiers) {
    auto tokens = getTokens("piñata + caffè * élite", KEYWORDS_ENGLISH);

    ASSERT_EQ(tokens.size(), 6); 
    AssertToken(tokens[0], TokenType::IDENTIFIER,  "piñata");
    AssertToken(tokens[1], TokenType::PLUS,        "+");
    AssertToken(tokens[2], TokenType::IDENTIFIER,  "caffè");
    AssertToken(tokens[3], TokenType::MULTIPLY,    "*");
    AssertToken(tokens[4], TokenType::IDENTIFIER,  "élite");
    AssertToken(tokens[5], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, HandlesComplexIdentifiers) {
    auto tokens = getTokens("_var123_tempé", KEYWORDS_ENGLISH);
    
    ASSERT_EQ(tokens.size(), 2);
    AssertToken(tokens[0], TokenType::IDENTIFIER,  "_var123_tempé");
    AssertToken(tokens[1], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, IgnoresExtraneousWhitespace) {
    auto tokens = getTokens("  89   \n \t  +  \r  variable  ", KEYWORDS_ENGLISH);
    
    ASSERT_EQ(tokens.size(), 4);
    AssertToken(tokens[0], TokenType::NUMBER,      "89");
    AssertToken(tokens[1], TokenType::PLUS,        "+");
    AssertToken(tokens[2], TokenType::IDENTIFIER,  "variable");
    AssertToken(tokens[3], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, HandlesCompoundOperators) {
    auto tokens = getTokens("x += 5 == y && z || a *= 2", KEYWORDS_ENGLISH);

    ASSERT_EQ(tokens.size(), 12);
    AssertToken(tokens[0], TokenType::IDENTIFIER,   "x");
    AssertToken(tokens[1], TokenType::PLUS_EQUALS,  "+=");
    AssertToken(tokens[2], TokenType::NUMBER,       "5");
    AssertToken(tokens[3], TokenType::EQUALS,       "==");
    AssertToken(tokens[4], TokenType::IDENTIFIER,   "y");
    AssertToken(tokens[5], TokenType::AND,          "&&");
    AssertToken(tokens[6], TokenType::IDENTIFIER,   "z");
    AssertToken(tokens[7], TokenType::OR,           "||");
    AssertToken(tokens[8], TokenType::IDENTIFIER,   "a");
    AssertToken(tokens[9], TokenType::STAR_EQUALS,  "*=");
    AssertToken(tokens[10], TokenType::NUMBER,      "2");
    AssertToken(tokens[11], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, DistinguishesSingleFromDouble) {
    auto tokens = getTokens("+ += = ==", KEYWORDS_ENGLISH);
    
    ASSERT_EQ(tokens.size(), 5);
    AssertToken(tokens[0], TokenType::PLUS,        "+");
    AssertToken(tokens[1], TokenType::PLUS_EQUALS, "+=");
    AssertToken(tokens[2], TokenType::ASSIGNMENT,  "=");
    AssertToken(tokens[3], TokenType::EQUALS,      "==");
    AssertToken(tokens[4], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, EmptyString) {
    auto tokens = getTokens("", KEYWORDS_ENGLISH);
    
    ASSERT_EQ(tokens.size(), 1);
    AssertToken(tokens[0], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, BlankString) {
    auto tokens = getTokens(" ", KEYWORDS_ENGLISH);
    
    ASSERT_EQ(tokens.size(), 1);
    AssertToken(tokens[0], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, AllDialectKeywordsAreEqual) {
    ASSERT_EQ(KEYWORDS_ENGLISH.size(), 24);
    ASSERT_EQ(KEYWORDS_FRENCH.size(),  24);
}

TEST_F(LexerTest, FrencDialect) {
    auto tokens = getTokens("dialecte constante fonction sinon si tantque demarrer afficher retourner structure variable entier entier64 octet erreur decimal chaine booleen et ou vrai faux variadique zero", KEYWORDS_FRENCH);

    ASSERT_EQ(tokens.size(), 25);
    AssertToken(tokens[0],  TokenType::DIALECT,     "dialecte");
    AssertToken(tokens[1],  TokenType::CONST,       "constante");
    AssertToken(tokens[3],  TokenType::ELSE,        "sinon");
    AssertToken(tokens[4],  TokenType::IF,          "si");
    AssertToken(tokens[5],  TokenType::WHILE,       "tantque");
    AssertToken(tokens[6],  TokenType::MAIN,        "demarrer");
    AssertToken(tokens[7],  TokenType::PRINT,       "afficher");
    AssertToken(tokens[8],  TokenType::RETURN,      "retourner");
    AssertToken(tokens[9],  TokenType::STRUCT,      "structure");
    AssertToken(tokens[10], TokenType::LET,         "variable");
    AssertToken(tokens[11], TokenType::INT,         "entier");
    AssertToken(tokens[12], TokenType::INT,         "entier64");
    AssertToken(tokens[13], TokenType::BYTE,        "octet");
    AssertToken(tokens[14], TokenType::ERROR,       "erreur");
    AssertToken(tokens[15], TokenType::FLOAT,       "decimal");
    AssertToken(tokens[16], TokenType::STRING,      "chaine");
    AssertToken(tokens[17], TokenType::BOOL,        "booleen");
    AssertToken(tokens[18], TokenType::AND,         "et");
    AssertToken(tokens[19], TokenType::OR,          "ou");
    AssertToken(tokens[20], TokenType::TRUE,        "vrai");
    AssertToken(tokens[21], TokenType::FALSE,       "faux");
    AssertToken(tokens[22], TokenType::VARIADIC,    "variadique");
    AssertToken(tokens[23], TokenType::ZERO,        "zero");
    AssertToken(tokens[24], TokenType::END_OF_FILE, "");
}

int main (int argc, char *argv[]) {
  testing::InitGoogleTest();
  int _ = RUN_ALL_TESTS();
  return 0;
}

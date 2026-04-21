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

void AssertToken(const Token& t, TokenType expectedType, const std::string& expectedValue) {
    EXPECT_EQ(t.type, expectedType);
    EXPECT_EQ(t.value, expectedValue);
}

TEST_F(LexerTest, HandlesBasicArithmetic) {
  auto tokens = getTokens("123 + 456", KEYWORDS_ENGLISH);

  ASSERT_EQ(tokens.size(), 4); 
  AssertToken(tokens[0], TokenType::NUMBER, "123");
  AssertToken(tokens[1], TokenType::PLUS, "+");
  AssertToken(tokens[2], TokenType::NUMBER, "456");
  AssertToken(tokens[3], TokenType::END_OF_FILE, "");
} 

TEST_F(LexerTest, HandlesAccentedIdentifiers) {
    auto tokens = getTokens("piñata + caffè * élite", KEYWORDS_ENGLISH);

    ASSERT_EQ(tokens.size(), 6); 
    AssertToken(tokens[0], TokenType::IDENTIFIER, "piñata");
    AssertToken(tokens[1], TokenType::PLUS, "+");
    AssertToken(tokens[2], TokenType::IDENTIFIER, "caffè");
    AssertToken(tokens[3], TokenType::MULTIPLY, "*");
    AssertToken(tokens[4], TokenType::IDENTIFIER, "élite");
    AssertToken(tokens[5], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, HandlesComplexIdentifiers) {
    auto tokens = getTokens("_var123_tempé", KEYWORDS_ENGLISH);
    
    ASSERT_EQ(tokens.size(), 2);
    AssertToken(tokens[0], TokenType::IDENTIFIER, "_var123_tempé");
    AssertToken(tokens[1], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, IgnoresExtraneousWhitespace) {
    auto tokens = getTokens("  89   \n \t  +  \r  variable  ", KEYWORDS_ENGLISH);
    
    ASSERT_EQ(tokens.size(), 4);
    AssertToken(tokens[0], TokenType::NUMBER, "89");
    AssertToken(tokens[1], TokenType::PLUS, "+");
    AssertToken(tokens[2], TokenType::IDENTIFIER, "variable");
    AssertToken(tokens[3], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, HandlesCompoundOperators) {
    auto tokens = getTokens("x += 5 == y && z || a *= 2", KEYWORDS_ENGLISH);

    ASSERT_EQ(tokens.size(), 12);
    AssertToken(tokens[0], TokenType::IDENTIFIER, "x");
    AssertToken(tokens[1], TokenType::PLUS_EQUALS, "+=");
    AssertToken(tokens[2], TokenType::NUMBER, "5");
    AssertToken(tokens[3], TokenType::EQUALS, "==");
    AssertToken(tokens[4], TokenType::IDENTIFIER, "y");
    AssertToken(tokens[5], TokenType::AND, "&&");
    AssertToken(tokens[6], TokenType::IDENTIFIER, "z");
    AssertToken(tokens[7], TokenType::OR, "||");
    AssertToken(tokens[8], TokenType::IDENTIFIER, "a");
    AssertToken(tokens[9], TokenType::STAR_EQUALS, "*=");
    AssertToken(tokens[10], TokenType::NUMBER, "2");
    AssertToken(tokens[11], TokenType::END_OF_FILE, "");
}

TEST_F(LexerTest, DistinguishesSingleFromDouble) {
    auto tokens = getTokens("+ += = ==", KEYWORDS_ENGLISH);
    
    ASSERT_EQ(tokens.size(), 5);
    AssertToken(tokens[0], TokenType::PLUS, "+");
    AssertToken(tokens[1], TokenType::PLUS_EQUALS, "+=");
    AssertToken(tokens[2], TokenType::ASSIGNMENT, "=");
    AssertToken(tokens[3], TokenType::EQUALS, "==");
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

int main (int argc, char *argv[]) {
  testing::InitGoogleTest();
  int _ = RUN_ALL_TESTS();
  return 0;
}

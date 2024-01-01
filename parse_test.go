package parse

import (
	"testing"
)

func TestTokenization(t *testing.T) {
	input := "Benchpress @90kg 5*5 \n Benchpress @87,5lbs 5/5 - \r \t $"
	expectedTokens := []Token{
		{String, "Benchpress", "Benchpress", 1},
		{WhiteSpace, " ", nil, 1},
		{Asperand, "@", nil, 1},
		{Number, "90", 90.0, 1},
		{WeightUnit, "kg", "kg", 1},
		{WhiteSpace, " ", nil, 1},
		{Number, "5", 5.0, 1},
		{Asterisk, "*", nil, 1},
		{Number, "5", 5.0, 1},
		{WhiteSpace, " ", nil, 1},
		{Newline, "\n", nil, 1},
		{WhiteSpace, " ", nil, 2},
		{String, "Benchpress", "Benchpress", 2},
		{WhiteSpace, " ", nil, 2},
		{Asperand, "@", nil, 2},
		{Number, "87.5", 87.5, 2},
		{WeightUnit, "lbs", "lbs", 2},
		{WhiteSpace, " ", nil, 2},
		{Number, "5", 5.0, 2},
		{ForwardSlash, "/", nil, 2},
		{Number, "5", 5.0, 2},
		{WhiteSpace, " ", nil, 2},
		{Hyphen, "-", nil, 2},
		{WhiteSpace, " ", nil, 2},
		{WhiteSpace, "\r", nil, 2},
		{WhiteSpace, " ", nil, 2},
		{WhiteSpace, "\t", nil, 2},
		{WhiteSpace, " ", nil, 2},
		{EOF, "", nil, 2},
	}

	s := newScanner([]rune(input))
	tokens, errs := s.scan()
	err := errs[0]

	if err.Error() != "unexpected character \"$\" at line 2" {
		t.Errorf("Expected error %q, got %q", "unexpected character \"$\" at line 2", err.Error())
	}

	for i, token := range tokens {
		if token.variant != expectedTokens[i].variant {
			t.Errorf("Expected token variant %v, got %v", expectedTokens[i].variant, token.variant)
		}
		if token.lexeme != expectedTokens[i].lexeme {
			t.Errorf("Expected lexeme %v, got %v", expectedTokens[i].lexeme, token.lexeme)
		}
		if token.literal != expectedTokens[i].literal {
			t.Errorf("Expected literal %v, got %v", expectedTokens[i].literal, token.literal)
		}
		if token.line != expectedTokens[i].line {
			t.Errorf("Expected token in line %v, got line %v", expectedTokens[i].line, token.line)
		}
	}

}

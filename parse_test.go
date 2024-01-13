package main

import (
	"testing"

	. "github.com/koeq/gl-parser/scanner"
	. "github.com/koeq/gl-parser/types"
)

func TestTokenization(t *testing.T) {
	input := "Benchpress @90kg 5*5 \n Benchpress @87,5lbs 5/5 - \r \t $"
	expectedTokens := []Token{
		{Variant: String, Lexeme: "Benchpress", Literal: "Benchpress", Line: 1},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 1},
		{Variant: Asperand, Lexeme: "@", Literal: nil, Line: 1},
		{Variant: Number, Lexeme: "90", Literal: 90.0, Line: 1},
		{Variant: WeightUnit, Lexeme: "kg", Literal: "kg", Line: 1},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 1},
		{Variant: Number, Lexeme: "5", Literal: 5.0, Line: 1},
		{Variant: Asterisk, Lexeme: "*", Literal: nil, Line: 1},
		{Variant: Number, Lexeme: "5", Literal: 5.0, Line: 1},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 1},
		{Variant: Newline, Lexeme: "\n", Literal: nil, Line: 1},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 2},
		{Variant: String, Lexeme: "Benchpress", Literal: "Benchpress", Line: 2},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 2},
		{Variant: Asperand, Lexeme: "@", Literal: nil, Line: 2},
		{Variant: Number, Lexeme: "87.5", Literal: 87.5, Line: 2},
		{Variant: WeightUnit, Lexeme: "lbs", Literal: "lbs", Line: 2},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 2},
		{Variant: Number, Lexeme: "5", Literal: 5.0, Line: 2},
		{Variant: ForwardSlash, Lexeme: "/", Literal: nil, Line: 2},
		{Variant: Number, Lexeme: "5", Literal: 5.0, Line: 2},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 2},
		{Variant: Hyphen, Lexeme: "-", Literal: nil, Line: 2},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 2},
		{Variant: WhiteSpace, Lexeme: "\r", Literal: nil, Line: 2},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 2},
		{Variant: WhiteSpace, Lexeme: "\t", Literal: nil, Line: 2},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 2},
		{Variant: EOF, Lexeme: "", Literal: nil, Line: 2},
	}

	sc := NewScanner([]rune(input))
	tokens, errs := sc.Scan()
	err := errs[0]

	if err.Error() != "unexpected character \"$\" at line 2" {
		t.Errorf("Expected error %q, got %q", "unexpected character \"$\" at line 2", err.Error())
	}

	for i, token := range tokens {
		if token.Variant != expectedTokens[i].Variant {
			t.Errorf("Expected token variant %v, got %v", expectedTokens[i].Variant, token.Variant)
		}
		if token.Lexeme != expectedTokens[i].Lexeme {
			t.Errorf("Expected lexeme %v, got %v", expectedTokens[i].Lexeme, token.Lexeme)
		}
		if token.Literal != expectedTokens[i].Literal {
			t.Errorf("Expected literal %v, got %v", expectedTokens[i].Literal, token.Literal)
		}
		if token.Line != expectedTokens[i].Line {
			t.Errorf("Expected token in line %v, got line %v", expectedTokens[i].Line, token.Line)
		}
	}

}

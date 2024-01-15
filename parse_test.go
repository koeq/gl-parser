package parser

import (
	"reflect"
	"testing"

	. "github.com/koeq/gl-parser/scanner"
	. "github.com/koeq/gl-parser/types"
)

func TestTokenization(t *testing.T) {
	input := "Bench Press @90kg 5*5 \n Bench Press @87,5lbs 5/5 - \r \t $"
	expectedTokens := []Token{
		{Variant: String, Lexeme: "Bench", Literal: "Bench", Line: 1},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 1},
		{Variant: String, Lexeme: "Press", Literal: "Press", Line: 1},
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
		{Variant: String, Lexeme: "Bench", Literal: "Bench", Line: 2},
		{Variant: WhiteSpace, Lexeme: " ", Literal: nil, Line: 2},
		{Variant: String, Lexeme: "Press", Literal: "Press", Line: 2},
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
			t.Errorf("Expected variant %v, got %v", expectedTokens[i].Variant, token.Variant)
		}
		if token.Lexeme != expectedTokens[i].Lexeme {
			t.Errorf("Expected %v, got %v", expectedTokens[i].Lexeme, token.Lexeme)
		}
		if token.Literal != expectedTokens[i].Literal {
			t.Errorf("Expected %v, got %v", expectedTokens[i].Literal, token.Literal)
		}
		if token.Line != expectedTokens[i].Line {
			t.Errorf("Expected token in line %v, got line %v", expectedTokens[i].Line, token.Line)
		}
	}
}

func TestParsing(t *testing.T) {
	input := "Bench Press @90kg 5/5/5 \n Squats @100kg 3*10 @140lbs 10 \n Dumbbell#-Rows @20 8/8/"

	expected := []Exercise{
		{
			Name: "Bench Press",
			Weight: Weight{
				Value: 90.0,
				Unit:  "kg",
			},
			Reps: []int{5, 5, 5},
		},
		{
			Name: "Squats",
			Weight: Weight{
				Value: 100.0,
				Unit:  "kg",
			},
			Reps: []int{10, 10, 10},
		},
		{
			Name: "Squats",
			Weight: Weight{
				Value: 140.0,
				Unit:  "lbs",
			},
			Reps: []int{10},
		},
		{
			Name: "Dumbbell-Rows",
			Weight: Weight{
				Value: 20.0,
				Unit:  "kg",
			},
			Reps: []int{8, 8},
		},
	}

	if exercises, _ := Parse(input); reflect.DeepEqual(exercises, expected) != true {
		t.Errorf("Expected exercises %v, got %v", expected, exercises)
	}
}

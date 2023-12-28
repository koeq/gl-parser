package parse

import (
	"errors"
	"fmt"
)

// Weight Units
const (
	Kilogram = "kg"
	Pound    = "lbs"
)

type TokenVariant string

// Token Variants
const (
	Asperand     TokenVariant = "ASPERAND"
	Asterisk     TokenVariant = "ASTERISK"
	EOF          TokenVariant = "EOF"
	ForwardSlash TokenVariant = "FORWARD_SLASH"
	Hyphen       TokenVariant = "HYPHEN"
	Newline      TokenVariant = "NEWLINE"
	Number       TokenVariant = "NUMBER"
	String       TokenVariant = "STRING"
	WeightUnit   TokenVariant = "WEIGHT_UNIT"
	Whitespace   TokenVariant = "WHITE_SPACE"
)

var tokenMap = map[rune]TokenVariant{
	'@':  Asperand,
	'*':  Asterisk,
	'/':  ForwardSlash,
	'-':  Hyphen,
	'\n': Newline,
	' ':  Whitespace,
	'\r': Whitespace,
	'\t': Whitespace,
}

type Weight struct {
	value float64
	unit  string
}

type Exercise struct {
	name   string
	weight Weight
	reps   []int
}

type Token struct {
	variant TokenVariant
	lexeme  string
	line    int
}

type Scanner struct {
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *Scanner) isAtEnd(src []rune) bool {
	return s.current >= len(src)
}

func (s *Scanner) advance(src []rune) rune {
	next := src[s.current]
	s.current++

	return next
}

func (s *Scanner) scan(src []rune) error {
	r := s.advance(src)

	switch r {
	case '@', '*', '/', '-', '\n', ' ', '\r', '\t':
		addToken(tokenMap[r])
	default:
		switch r {
		case isString(r):
		case isNum(r):
		default:
			return fmt.Errorf("unexpected character at line %d", s.line)
		}

	}
}

func (s *Scanner) tokenize(src []rune) (tokens []Token) {
	for !s.isAtEnd(src) {
		s.start = s.current
		s.scan(src)
	}

	tokens = append(s.tokens, Token{variant: "EOF", lexeme: "", line: s.line})

	return s.tokens
}

type Interpreter struct{}

func (i *Interpreter) interpret(tokens []Token) (exercises []Exercise, err error) {

	return []Exercise{}, nil
}

func Parse(source string) (exercises []Exercise, err error) {
	if source == "" {
		return nil, errors.New("empty source string")
	}

	tokens := (&Scanner{}).tokenize([]rune(source))

	return (&Interpreter{}).interpret(tokens)
}

// Utils
func isValidTokenVariant(tv TokenVariant) bool {
	switch tv {
	case Asperand, ForwardSlash, Asterisk, Number, String, Hyphen, WeightUnit, Whitespace:
		return true
	default:
		return false
	}
}

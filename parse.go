package parse

import (
	"errors"
	"fmt"
	"regexp"
)

type TokenVariant string

// Token Variants
// The token type String represents one or multiple letters that are not keywords
const (
	Asperand     TokenVariant = "ASPERAND"
	Asterisk     TokenVariant = "ASTERISK"
	EOF          TokenVariant = "EOF"
	ForwardSlash TokenVariant = "FORWARD_SLASH"
	Hyphen       TokenVariant = "HYPHEN"
	Newline      TokenVariant = "NEWLINE"
	Digit        TokenVariant = "DIGIT"
	String       TokenVariant = "STRING"
	WeightUnit   TokenVariant = "WEIGHT_UNIT"
	WhiteSpace   TokenVariant = "WHITE_SPACE"
)

var tokenVariantMap = map[rune]TokenVariant{
	'@':  Asperand,
	'*':  Asterisk,
	'/':  ForwardSlash,
	'-':  Hyphen,
	'\n': Newline,
	' ':  WhiteSpace,
	'\r': WhiteSpace,
	'\t': WhiteSpace,
}

var keywordVariantMap = map[string]TokenVariant{
	"kg":  WeightUnit,
	"lbs": WeightUnit,
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
	src     []rune
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.src)
}

func (s *Scanner) advance() rune {
	next := s.src[s.current]
	s.current++

	return next
}

func (s *Scanner) addToken(variant TokenVariant, literal string) {
	s.tokens = append(s.tokens, Token{variant, literal, s.line})
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	return s.src[s.current]
}

func (s *Scanner) processWord() {
	for isLetter(s.peek()) {
		s.advance()
	}

	word := string(s.src[s.start:s.current])
	tokenVariant, isKeyword := keywordVariantMap[word]

	if isKeyword {
		s.addToken(tokenVariant, word)
	} else {
		s.addToken(String, word)
	}
}

func (s *Scanner) scan() error {
	r := s.advance()

	switch r {
	case '@', '*', '/', '-', '\n', ' ', '\r', '\t':
		s.addToken(tokenVariantMap[r], "")
	default:
		switch r {
		case isLetter(r):
		case isDigit(r):
		default:
			return fmt.Errorf("unexpected character at line %d", s.line)
		}
	}

	return nil
}

func (s *Scanner) tokenize() (tokens []Token) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scan()
	}

	s.tokens = append(s.tokens, Token{variant: "EOF", lexeme: "", line: s.line})

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

	// TODO: create scanner
	tokens := (&Scanner{}).tokenize()

	// TODO: create interpreter
	return (&Interpreter{}).interpret(tokens)
}

// Utils
func isValidTokenVariant(tv TokenVariant) bool {
	switch tv {
	case Asperand, Asterisk, Digit, ForwardSlash, Hyphen, String, WeightUnit, WhiteSpace:
		return true
	default:
		return false
	}
}

var (
	// match letters + combining marks
	letterRegex = regexp.MustCompile(`[\p{L}\p{M}]`)
	digitRegex  = regexp.MustCompile(`\d`)
)

func isLetter(r rune) bool {
	return letterRegex.MatchString(string(r))
}

func isDigit(r rune) bool {
	return digitRegex.MatchString(string(r))
}

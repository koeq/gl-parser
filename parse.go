package parse

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TokenVariant string

// token variants
// "STRING" represents one or multiple letters that are not keywords
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
	WhiteSpace   TokenVariant = "WHITE_SPACE"
)

var TokenVariantMap = map[rune]TokenVariant{
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
	literal interface{}
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

func (s *Scanner) isNextAtEnd() bool {
	return s.current+1 >= len(s.src)
}

func (s *Scanner) advance() rune {
	next := s.src[s.current]
	s.current++

	return next
}

func (s *Scanner) addToken(variant TokenVariant, lexeme string, literal interface{}) {
	s.tokens = append(s.tokens, Token{variant, lexeme, literal, s.line})

	if variant == Newline {
		s.line++
	}
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
		s.addToken(tokenVariant, word, word)
	} else {
		s.addToken(String, word, word)
	}
}

func (s *Scanner) peekNext() rune {
	if s.isNextAtEnd() {
		return 0
	}

	return s.src[s.current+1]
}

func (s *Scanner) processNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// look for fractional part
	if s.peek() == '.' || s.peek() == ',' && isDigit(s.peekNext()) {
		// consume it
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	sNum := strings.Replace(string(s.src[s.start:s.current]), ",", ".", 1)
	f, err := strconv.ParseFloat(sNum, 32)

	// TODO: handle error case
	if err == nil {
		s.addToken(Number, sNum, f)
	}
}

func (s *Scanner) tokenize() error {
	r := s.advance()

	switch r {
	case '@', '*', '/', '-', '\n', ' ', '\r', '\t':
		s.addToken(TokenVariantMap[r], string(r), nil)
	default:
		switch {
		case isLetter(r):
			s.processWord()
		case isDigit(r):
			s.processNumber()
		default:
			return fmt.Errorf("unexpected character at line %d", s.line)
		}
	}

	return nil
}

func (s *Scanner) scan() (tokens []Token) {
	for !s.isAtEnd() {
		s.start = s.current
		s.tokenize()
	}

	return append(s.tokens, Token{"EOF", "", nil, s.line})

}

type Interpreter struct{}

func (i *Interpreter) interpret(tokens []Token) (exercises []Exercise, err error) {

	return []Exercise{}, nil
}

func newScanner(src []rune) (s *Scanner) {
	return &Scanner{src: src, tokens: []Token{}, start: 0, current: 0, line: 1}
}

func Parse(source string) (exercises []Exercise, err error) {
	if source == "" {
		return nil, errors.New("empty source string")
	}

	scanner := newScanner([]rune(source))
	tokens := scanner.scan()

	// TODO: create interpreter
	return (&Interpreter{}).interpret(tokens)
}

// utils
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

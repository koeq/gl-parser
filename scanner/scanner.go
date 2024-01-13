package scanner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	. "github.com/koeq/gl-parser/types"
)

type Scanner struct {
	src     []rune
	tokens  []Token
	start   int
	current int
	line    int
	errors  []ScanError
}

type ScanError struct {
	line int
	r    string
}

var tokenVariantMap = map[string]TokenVariant{
	"@":  Asperand,
	"*":  Asterisk,
	"/":  ForwardSlash,
	"-":  Hyphen,
	"\n": Newline,
	" ":  WhiteSpace,
	"\r": WhiteSpace,
	"\t": WhiteSpace,
}

var keywordVariantMap = map[string]TokenVariant{
	"kg":  WeightUnit,
	"lbs": WeightUnit,
}

var wordRegex = regexp.MustCompile(`[\p{L}\p{M}]+`) // match letters + combining marks
var numberRegex = regexp.MustCompile(`\d+(\.\d+)?`) // match number with optional fractional part

func (e *ScanError) Error() string {
	return fmt.Sprintf("unexpected character %q at line %d", e.r, e.line)
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.src)
}

func (sc *Scanner) isNextAtEnd() bool {
	return sc.current+1 >= len(sc.src)
}

func (sc *Scanner) advance() string {
	next := string(sc.src[sc.current])
	sc.current++

	return next
}

func (sc *Scanner) addToken(variant TokenVariant, lexeme string, literal interface{}) {
	sc.tokens = append(sc.tokens, Token{Variant: variant, Lexeme: lexeme, Literal: literal, Line: sc.line})

	if variant == Newline {
		sc.line++
	}
}

func (sc *Scanner) peek() string {
	if sc.isAtEnd() {
		return ""
	}

	return string(sc.src[sc.current])
}

func isWord(s string) bool {
	return wordRegex.MatchString(s)
}

func (sc *Scanner) processWord() {
	for isWord(sc.peek()) {
		sc.advance()
	}

	word := string(sc.src[sc.start:sc.current])
	tokenVariant, isKeyword := keywordVariantMap[word]

	if isKeyword {
		sc.addToken(tokenVariant, word, word)
	} else {
		sc.addToken(String, word, word)
	}
}

func (sc *Scanner) peekNext() string {
	if sc.isNextAtEnd() {
		return ""
	}

	return string(sc.src[sc.current+1])
}

func isNumber(s string) bool {
	return numberRegex.MatchString(s)
}

func (sc *Scanner) processNumber() {
	for isNumber(sc.peek()) {
		sc.advance()
	}

	// look for fractional part
	if sc.peek() == "." || sc.peek() == "," && isNumber(sc.peekNext()) {
		// consume it
		sc.advance()

		for isNumber(sc.peek()) {
			sc.advance()
		}
	}

	sNum := strings.Replace(string(sc.src[sc.start:sc.current]), ",", ".", 1)
	f, err := strconv.ParseFloat(sNum, 32)

	if err != nil {
		sc.errors = append(sc.errors, ScanError{sc.line, sNum})
	} else {
		sc.addToken(Number, sNum, f)
	}
}

func (sc *Scanner) tokenize() {
	s := sc.advance()

	switch s {
	case "@", "*", "/", "-", "\n", " ", "\r", "\t":
		sc.addToken(tokenVariantMap[s], s, nil)
	default:
		switch {
		case isWord(s):
			sc.processWord()
		case isNumber(s):
			sc.processNumber()
		default:
			sc.errors = append(sc.errors, ScanError{sc.line, s})
		}
	}
}

func (sc *Scanner) Scan() (tokens []Token, errs []ScanError) {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.tokenize()
	}

	sc.addToken(EOF, "", nil)

	return sc.tokens, sc.errors
}

func NewScanner(src []rune) (sc *Scanner) {
	return &Scanner{src: src, tokens: []Token{}, start: 0, current: 0, line: 1}
}

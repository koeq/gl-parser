package main

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
	errors  []ScanError
}

type ScanError struct {
	line int
	r    string
}

type Exercise struct {
	name   string
	weight Weight
	reps   []int
}

type Weight struct {
	value float64
	unit  string
}

type Interpreter struct {
	tokens        []Token
	exercises     []Exercise
	start         int
	current       int
	exerciseIndex int
	exerciseName  string
	weight        Weight
	reps          []int
}

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
	sc.tokens = append(sc.tokens, Token{variant, lexeme, literal, sc.line})

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

func (sc *Scanner) scan() (tokens []Token, errs []ScanError) {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.tokenize()
	}

	sc.addToken(EOF, "", nil)

	return sc.tokens, sc.errors
}

func newScanner(src []rune) (sc *Scanner) {
	return &Scanner{src: src, tokens: []Token{}, start: 0, current: 0, line: 1}
}

func (it *Interpreter) isAtEnd() bool {
	return it.tokens[it.current].variant == EOF
}

func (it *Interpreter) advance() Token {
	token := it.tokens[it.current]
	it.current++

	return token
}

func (it *Interpreter) peek() Token {
	return it.tokens[it.current]
}

func (it *Interpreter) build() string {
	var builder strings.Builder
	ts := it.tokens[it.start:it.current]

	for _, t := range ts {
		builder.WriteString(strings.TrimSpace(t.lexeme))
	}

	return builder.String()
}

func (it *Interpreter) processExerciseName(token Token) {
	for isExerciseName(token.variant) {
		next := it.peek()

		if !isExerciseName(next.variant) {
			break
		}

		token = it.advance()
	}

	// reset interpreter state
	it.weight.unit = ""
	it.weight.value = 0
	it.reps = nil

	it.exerciseIndex++

	name := it.build()
	it.exercises = append(it.exercises, Exercise{name, it.weight, it.reps})
}

// TODO: index out of bounds error happening here
func (it *Interpreter) processWeight() {
	next := it.peek()

	if next.variant == Number {
		// TODO: failed assertion would cause a runtime panic -> find better solution or handle error
		it.weight.value = next.literal.(float64)
		it.advance()
	}

	next = it.peek()

	if next.variant == WeightUnit {
		// TODO: failed assertion causes a runtime panic -> find better solution or handle error
		it.weight.unit = next.literal.(string)
		it.advance()
	}

	// if there is already a weight we want create a second exercise with the same name
	// e.g. Benchpress @100 8/8 @95 8/8 -> Benchpress 100kg 8/8
	//                                  -> Benchpress 95kg 8/8
	if it.weight.value != 0 {
		it.exerciseIndex++
		it.reps = nil
	}

	it.exercises[it.exerciseIndex] = Exercise{it.exerciseName, it.weight, it.reps}
}

func (it *Interpreter) processReps(token Token) {
	for isReps(token.variant) {
		next := it.peek()

		if !isReps(next.variant) {
			break
		}

		token = it.advance()
	}

	repStr := it.build()

	// format: int*int
	if isRepsMultiplierFormat(repStr) {
		multiplierReps := strings.Split(repStr, "*")
		multiplier, mErr := strconv.Atoi(multiplierReps[0])
		repCount, rErr := strconv.Atoi(multiplierReps[1])

		if mErr != nil && rErr != nil {
			fmt.Println("error parsing reps: ", mErr, rErr)

			return
		}

		it.reps = make([]int, 0, multiplier)

		for i := 0; i < multiplier; i++ {
			it.reps = append(it.reps, repCount)
		}
	}

	// format: int/int/int
	if isRepsEnumerationFormat(repStr) {
		numStrs := intRegex.FindAllString(repStr, -1)
		it.reps = make([]int, 0, len(numStrs))

		for _, s := range numStrs {

			num, err := strconv.Atoi(s)

			if err != nil {
				fmt.Println("error parsing reps: ", err)

				continue
			}

			it.reps = append(it.reps, num)
		}
	}

	it.exercises[it.exerciseIndex] = Exercise{it.exerciseName, it.weight, it.reps}
}

func (it *Interpreter) interpret() (exercises []Exercise, err error) {
	for !it.isAtEnd() {
		it.start = it.current
		token := it.advance()

		switch token.variant {
		case "HYPHEN", "STRING":
			it.processExerciseName(token)

		case "ASPERAND":
			it.processWeight()

		case "NUMBER":
			it.processReps(token)
		}
		// TODO: what happens if a different token variant is encountered (e.g. asterisk)?

	}

	return it.exercises, nil
}

func newInterpreter(tokens []Token) (it *Interpreter) {
	return &Interpreter{tokens: tokens, exercises: []Exercise{}, start: 0, current: 0, exerciseIndex: 0, weight: Weight{0, ""}, reps: nil}
}

func Parse(source string) (exercises []Exercise, err error) {
	if source == "" {
		return nil, errors.New("empty source string")
	}

	sc := newScanner([]rune(source))
	tokens, errs := sc.scan()

	for _, err := range errs {
		fmt.Println(err)
	}

	it := newInterpreter(tokens)

	return it.interpret()
}

// utils
var (
	wordRegex                  = regexp.MustCompile(`[\p{L}\p{M}]+`) // match letters + combining marks
	numberRegex                = regexp.MustCompile(`\d+(\.\d+)?`)   // match number with optional fractional part
	repsMultiplierFormatRegex  = regexp.MustCompile(`\d+\*\d+`)      // match reps in format int*int
	repsEnumerationFormatRegex = regexp.MustCompile(`\d+(,\d+)+`)    // match reps in format int/int/int
	intRegex                   = regexp.MustCompile(`\d+`)           // match any int

)

func isWord(s string) bool {
	return wordRegex.MatchString(s)
}

func isNumber(s string) bool {
	return numberRegex.MatchString(s)
}

func isExerciseName(tv TokenVariant) bool {
	return tv == String || tv == Hyphen || tv == WhiteSpace
}

func isReps(tv TokenVariant) bool {
	return tv == Number || tv == ForwardSlash || tv == Asterisk
}

func isRepsMultiplierFormat(s string) bool {
	return repsMultiplierFormatRegex.MatchString(s)
}

func isRepsEnumerationFormat(s string) bool {
	return repsEnumerationFormatRegex.MatchString(s)
}

func main() {
	Parse("Benchpress @90kg")
}

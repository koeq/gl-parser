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
	tokens    []Token
	exercises []Exercise
	start     int
	current   int
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

func (in *Interpreter) isAtEnd() bool {
	return in.tokens[in.current].variant == EOF
}

func (in *Interpreter) advance() Token {
	token := in.tokens[in.current]
	in.current++

	return token
}

func (in *Interpreter) peek() Token {
	return in.tokens[in.current]
}

func (in *Interpreter) buildStr() string {
	var builder strings.Builder
	ts := in.tokens[in.start:in.current]

	for _, t := range ts {
		builder.WriteString(strings.TrimSpace(t.lexeme))
	}

	return builder.String()
}

func (in *Interpreter) consumeWhile(token Token, predicate func(tv TokenVariant) bool) Token {
	for predicate(token.variant) {
		next := in.peek()

		if !predicate(next.variant) {
			break
		}

		token = in.advance()
	}

	return token
}

func (in *Interpreter) processExerciseName(token Token) {
	in.consumeWhile(token, isExerciseName)
	name := in.buildStr()
	in.exercises = append(in.exercises, Exercise{name, Weight{value: 0, unit: ""}, nil})
}

func (in *Interpreter) processWeight() {
	var weight Weight
	currExercise := &in.exercises[len(in.exercises)-1]

	next := in.peek()

	if next.variant == Number {
		// TODO: failed assertion would cause a runtime panic -> find better solution or handle error
		weight.value = next.literal.(float64)
		in.advance()
	}

	next = in.peek()

	if next.variant == WeightUnit {
		// TODO: failed assertion causes a runtime panic -> find better solution or handle error
		weight.unit = next.literal.(string)
		in.advance()
	} else {
		// TODO: provide config to specify weight unit
		// default weight unit
		weight.unit = "kg"
	}

	// if there is already a weight we want create a second exercise with the same name
	// --> 	Benchpress @100 8/8 @102 8/8  -> 	Benchpress 100kg 8/8
	//                                  	 		Benchpress 102kg 8/8
	if currExercise.weight.value != 0 {
		in.exercises = append(in.exercises, Exercise{currExercise.name, weight, nil})
	} else {
		currExercise.weight = weight
	}
}

func (in *Interpreter) processReps(token Token) {
	in.consumeWhile(token, isReps)
	repStr := in.buildStr()

	var reps []int
	var err error

	// int*int
	if isRepsMultiplierFormat(repStr) {
		reps, err = parseMultiplierFormat(repStr)
		// int/int/int
	} else if isRepsEnumerationFormat(repStr) {
		reps, err = parseEnumerationFormat(repStr)
	}

	if err != nil {
		// TODO: add interpreter error handling
	}

	currExercise := &in.exercises[len(in.exercises)-1]
	currExercise.reps = reps
}

func (in *Interpreter) interpret() (exercises []Exercise, err error) {
	for !in.isAtEnd() {
		in.start = in.current
		token := in.advance()

		switch token.variant {
		case "HYPHEN", "STRING":
			in.processExerciseName(token)

		case "ASPERAND":
			in.processWeight()

		case "NUMBER":
			in.processReps(token)
		default:
			// TODO: add error handling if unexpected token is encountered
		}
	}

	return in.exercises, nil
}

func newInterpreter(tokens []Token) (in *Interpreter) {
	return &Interpreter{tokens: tokens, exercises: []Exercise{}, start: 0, current: 0}
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

	in := newInterpreter(tokens)

	return in.interpret()
}

// utils
var (
	wordRegex                  = regexp.MustCompile(`[\p{L}\p{M}]+`)  // match letters + combining marks
	numberRegex                = regexp.MustCompile(`\d+(\.\d+)?`)    // match number with optional fractional part
	repsMultiplierFormatRegex  = regexp.MustCompile(`\d+\*\d+`)       // match reps in format int*int
	repsEnumerationFormatRegex = regexp.MustCompile(`\d+(\/\d+)*\/?`) // match reps in format int/int/int
	intRegex                   = regexp.MustCompile(`\d+`)            // match any int

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

func parseMultiplierFormat(s string) ([]int, error) {
	var reps []int

	multiplierReps := strings.Split(s, "*")
	multiplier, err := strconv.Atoi(multiplierReps[0])

	if err != nil {
		return nil, err
	}

	repCount, err := strconv.Atoi(multiplierReps[1])

	if err != nil {
		return nil, err
	}

	reps = make([]int, 0, multiplier)

	for i := 0; i < multiplier; i++ {
		reps = append(reps, repCount)
	}

	return reps, nil
}

func parseEnumerationFormat(s string) ([]int, error) {
	numStrs := intRegex.FindAllString(s, -1)
	reps := make([]int, 0, len(numStrs))

	for _, s := range numStrs {

		num, err := strconv.Atoi(s)

		if err != nil {
			return nil, err
		}

		reps = append(reps, num)
	}

	return reps, nil
}

func main() {
	exercises, _ := Parse("Benchpress @90kg 5/5/5 \n Squats @100kg 5*10")
	fmt.Println(exercises)
}

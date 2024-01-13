package interpreter

import (
	"regexp"
	"strconv"
	"strings"

	. "github.com/koeq/gl-parser/types"
)

type Interpreter struct {
	tokens    []Token
	exercises []Exercise
	start     int
	current   int
}

var repsMultiplierFormatRegex = regexp.MustCompile(`\d+\*\d+`)        // match reps in format int*int
var repsEnumerationFormatRegex = regexp.MustCompile(`\d+(\/\d+)*\/?`) // match reps in format int/int/int
var intRegex = regexp.MustCompile(`\d+`)                              // match any int

func (in *Interpreter) isAtEnd() bool {
	return in.tokens[in.current].Variant == EOF
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
		builder.WriteString(strings.TrimSpace(t.Lexeme))
	}

	return builder.String()
}

func (in *Interpreter) consumeWhile(token Token, predicate func(tv TokenVariant) bool) Token {
	for predicate(token.Variant) {
		next := in.peek()

		if !predicate(next.Variant) {
			break
		}

		token = in.advance()
	}

	return token
}

func isExerciseName(tv TokenVariant) bool {
	return tv == String || tv == Hyphen || tv == WhiteSpace
}

func (in *Interpreter) processExerciseName(token Token) {
	in.consumeWhile(token, isExerciseName)
	name := in.buildStr()
	in.exercises = append(in.exercises, Exercise{Name: name, Weight: Weight{Value: 0, Unit: ""}, Reps: nil})
}

func (in *Interpreter) processWeight() {
	var weight Weight
	currExercise := &in.exercises[len(in.exercises)-1]

	next := in.peek()

	if next.Variant == Number {
		// TODO: failed assertion would cause a runtime panic -> find better solution or handle error
		weight.Value = next.Literal.(float64)
		in.advance()
	}

	next = in.peek()

	if next.Variant == WeightUnit {
		// TODO: failed assertion causes a runtime panic -> find better solution or handle error
		weight.Unit = next.Literal.(string)
		in.advance()
	} else {
		// TODO: provide config to specify weight unit
		// default weight unit
		weight.Unit = "kg"
	}

	// if there is already a weight we want create a second exercise with the same name
	// --> 	Benchpress @100 8/8 @102 8/8  -> 	Benchpress 100kg 8/8
	//                                  	 		Benchpress 102kg 8/8
	if currExercise.Weight.Value != 0 {
		in.exercises = append(in.exercises, Exercise{Name: currExercise.Name, Weight: weight, Reps: nil})
	} else {
		currExercise.Weight = weight
	}
}

func isReps(tv TokenVariant) bool {
	return tv == Number || tv == ForwardSlash || tv == Asterisk
}

func IsRepsMultiplierFormat(s string) bool {
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

func (in *Interpreter) processReps(token Token) {
	in.consumeWhile(token, isReps)
	repStr := in.buildStr()

	var reps []int
	var err error

	// int*int
	if IsRepsMultiplierFormat(repStr) {
		reps, err = parseMultiplierFormat(repStr)
		// int/int/int
	} else if isRepsEnumerationFormat(repStr) {
		reps, err = parseEnumerationFormat(repStr)
	}

	if err != nil {
		// TODO: add interpreter error handling
	}

	currExercise := &in.exercises[len(in.exercises)-1]
	currExercise.Reps = reps
}

func (in *Interpreter) Interpret() (exercises []Exercise, err error) {
	for !in.isAtEnd() {
		in.start = in.current
		token := in.advance()

		switch token.Variant {
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

func NewInterpreter(tokens []Token) (in *Interpreter) {
	return &Interpreter{tokens: tokens, exercises: []Exercise{}, start: 0, current: 0}
}

package interpreter

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	. "github.com/koeq/gl-parser/types"
)

type Interpreter struct {
	config    *Config
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
		builder.WriteString(t.Lexeme)
	}

	return builder.String()
}

func (in *Interpreter) advanceWhile(predicate func(tv TokenVariant) bool) {
	for predicate(in.peek().Variant) {
		in.advance()
	}

}

func isExerciseName(tv TokenVariant) bool {
	return tv == String || tv == Hyphen || tv == WhiteSpace
}

func (in *Interpreter) processExerciseName() {
	in.advanceWhile(isExerciseName)
	name := strings.TrimSpace(in.buildStr())
	in.exercises = append(in.exercises, Exercise{Name: name, Weight: Weight{Value: 0, Unit: ""}, Reps: nil})
}

func NewUnit(s string) (Unit, error) {
	if s == "kg" {
		return Metric, nil
	} else if s == "lbs" {
		return Imperial, nil
	}

	return "", errors.New("invalid weight unit")
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
		unit, err := NewUnit(next.Literal.(string))

		if err != nil {
			// TODO: handle err
		} else {
			weight.Unit = unit
		}

		in.advance()
	} else {
		// fallback to weight unit from config -> defaults to "kg"
		weight.Unit = in.config.WeightUnit
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

func (in *Interpreter) processReps() {
	in.advanceWhile(isReps)
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
			in.processExerciseName()

		case "ASPERAND":
			in.processWeight()

		case "NUMBER":
			in.processReps()
		default:
			// TODO: add error handling if unexpected token is encountered
		}
	}

	return in.exercises, nil
}

func NewInterpreter(tokens []Token, config *Config) (in *Interpreter) {
	return &Interpreter{config: config, tokens: tokens, exercises: []Exercise{}, start: 0, current: 0}
}

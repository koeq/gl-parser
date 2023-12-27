package parse

import (
	"errors"
)

const (
	Kilogram = "kg"
	Pound    = "lbs"
)

type Weight struct {
	value float64
	unit  string
}

type Exercise struct {
	name   string
	weight Weight
	reps   []int
}

type Token struct{}

type Scanner struct{}

type Interpreter struct{}

func (s *Scanner) Tokenize(source string) (tokens []Token) {

	return []Token{}
}

func (i *Interpreter) interpret(tokens []Token) (exercises []Exercise, err error) {

	return []Exercise{}, nil
}

func Parse(source string) (exercises []Exercise, err error) {
	if source == "" {
		return nil, errors.New("empty source string")
	}

	tokens := (&Scanner{}).Tokenize(source)

	return (&Interpreter{}).interpret(tokens)
}

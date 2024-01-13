package main

import (
	"errors"
	"fmt"

	. "github.com/koeq/gl-parser/interpreter"
	. "github.com/koeq/gl-parser/scanner"
	. "github.com/koeq/gl-parser/types"
)

func Parse(source string) (exercises []Exercise, err error) {
	if source == "" {
		return nil, errors.New("empty source string")
	}

	sc := NewScanner([]rune(source))
	tokens, errs := sc.Scan()

	for _, err := range errs {
		fmt.Println(err)
	}

	in := NewInterpreter(tokens)

	return in.Interpret()
}

func main() {
	exercises, _ := Parse("Benchpress @90kg 5/5/5 \n Squats @100kg 5*10")
	fmt.Println(exercises)
}

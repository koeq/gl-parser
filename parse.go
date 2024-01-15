package parser

import (
	"errors"
	"fmt"

	"github.com/koeq/gl-parser/interpreter"
	"github.com/koeq/gl-parser/scanner"
	. "github.com/koeq/gl-parser/types"
)

func Parse(source string) (exercises []Exercise, err error) {

	if source == "" {
		return nil, errors.New("empty source string")
	}

	sc := scanner.NewScanner([]rune(source))
	tokens, errs := sc.Scan()

	for _, err := range errs {
		fmt.Println(err.Error())
	}

	in := interpreter.NewInterpreter(tokens)

	return in.Interpret()
}

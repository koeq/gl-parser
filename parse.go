package parser

import (
	"errors"
	"fmt"

	"github.com/koeq/gl-parser/interpreter"
	"github.com/koeq/gl-parser/scanner"
	. "github.com/koeq/gl-parser/types"
)



func WithWeightUnit(unit Unit) ConfigOption {
	return func(c *Config) {
		if unit != Metric && unit != Imperial {
			panic(fmt.Sprintf("invalid weight unit: use '%s' or '%s' instead", Metric, Imperial))
		}

		c.WeightUnit = unit
	}
}

func NewConfig(opts ...ConfigOption) *Config {
	c := &Config{
		WeightUnit: Metric,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func Parse(source string, config *Config) (exercises []Exercise, err error) {

	if source == "" {
		return nil, errors.New("empty source string")
	}

	sc := scanner.NewScanner([]rune(source))
	tokens, errs := sc.Scan()

	for _, err := range errs {
		fmt.Println(err.Error())
	}

	in := interpreter.NewInterpreter(tokens, config)

	return in.Interpret()
}

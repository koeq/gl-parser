package main

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

func parse(source string) (exercises []Exercise, err error) {
	if source == "" {
		return nil, errors.New("empty source string")
	}

	return []Exercise{}, nil
}

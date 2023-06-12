package main

import (
	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func GenerateSign(input string) string {
	if len(input) > SignLen {
		input = input[:SignLen*3/4] + "..." + input[len(input)-SignLen/4-3:]
	}
	return input
}

package main

import "math/rand"

type Challenge struct {
	ElementA  int
	ElementB  int
	Answer    int
	Operation string
}

const MAX = 99
const MIN = 0

var ops = map[string]func(int, int) int{
	"+": func(a int, b int) int {
		return a + b
	},
}

func generateChallenge() Challenge {
	sum := "+"
	elementA := rand.Intn(MAX-MIN+1) + MIN
	elementB := rand.Intn(MAX-MIN+1) + MIN
	answer := ops[sum](elementA, elementB)

	return Challenge{
		elementA,
		elementB,
		answer,
		sum,
	}
}

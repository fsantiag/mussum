package main

import "math/rand"

// Challenge represents a math challenge
type Challenge struct {
	ElementA  int
	ElementB  int
	Answer    int
	Operation string
}

const max = 99
const min = 0

var ops = map[string]func(int, int) int{
	"+": func(a int, b int) int {
		return a + b
	},
}

func generateChallenge() Challenge {
	sum := "+"
	elementA := rand.Intn(max-min+1) + min
	elementB := rand.Intn(max-min+1) + min
	answer := ops[sum](elementA, elementB)

	return Challenge{
		elementA,
		elementB,
		answer,
		sum,
	}
}

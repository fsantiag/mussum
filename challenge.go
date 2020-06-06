package main

import (
	"math/rand"
	"time"
)

// Challenge represents a math challenge
type Challenge struct {
	ElementA  int
	ElementB  int
	Answer    int
	Operation string
}

const max = 100

var ops = map[string]func(int, int) int{
	"+": func(a int, b int) int {
		return a + b
	},
}

func generateChallenge() Challenge {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	sum := "+"
	elementA := r.Intn(max)
	elementB := r.Intn(max)
	answer := ops[sum](elementA, elementB)

	return Challenge{
		elementA,
		elementB,
		answer,
		sum,
	}
}

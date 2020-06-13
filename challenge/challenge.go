package challenge

import (
	"math/rand"
	"time"
)

// SumChallenge represents a math sum challenge
type SumChallenge struct {
	A         int
	B         int
	Result    int
	Operation string
}

const max = 100

// Generate will create a random sum challenge
func Generate() SumChallenge {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	A := r.Intn(max)
	B := r.Intn(max)
	result := A + B

	return SumChallenge{
		A,
		B,
		result,
		"+",
	}
}

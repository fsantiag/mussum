package challenge

import (
	"fmt"
	"math/rand"
	"time"
)

// SumChallenge represents a math sum challenge
type SumChallenge struct {
	ElementA  int
	ElementB  int
	Answer    int
	Operation string
}

func (c SumChallenge) String() string {
	return fmt.Sprintf("Qual o valor de %v + %v?", c.ElementA, c.ElementB)
}

const max = 100

// GenerateChallenge will create a random sum challenge
func Generate() SumChallenge {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	elementA := r.Intn(max)
	elementB := r.Intn(max)
	answer := elementA + elementB

	return SumChallenge{
		elementA,
		elementB,
		answer,
		"+",
	}
}

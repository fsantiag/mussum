package challenge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChallengeIsGenerated(t *testing.T) {
	challenge := Generate()

	answer := challenge.ElementA + challenge.ElementB

	assert.Equal(t, answer, challenge.Answer)
}

func TestRandomChallengesAreGenerated(t *testing.T) {
	c1 := Generate()
	c2 := Generate()

	assert.NotEqual(t, c1, c2)
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChallengeIsGenerated(t *testing.T) {
	challenge := generateChallenge()

	answer := challenge.ElementA + challenge.ElementB

	assert.Equal(t, answer, challenge.Answer)
}

func TestRandomChallengesAreGenerated(t *testing.T) {
	c1 := generateChallenge()
	c2 := generateChallenge()

	assert.NotEqual(t, c1, c2)
}

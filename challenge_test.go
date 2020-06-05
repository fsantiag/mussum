package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSumChallenge(t *testing.T) {
	challenge := generateChallenge()

	answer := challenge.ElementA + challenge.ElementB

	assert.Equal(t, answer, challenge.Answer)
}

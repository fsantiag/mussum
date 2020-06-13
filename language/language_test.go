package language

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultsToPortuguese(t *testing.T) {
	lang := GetDefault()

	assert.Equal(t, "pt", lang.ID())
}
func TestSelectLanguageFromEnvVar(t *testing.T) {
	os.Setenv("LANGUAGE", "en")

	lang := GetDefault()

	assert.Equal(t, "en", lang.ID())
}

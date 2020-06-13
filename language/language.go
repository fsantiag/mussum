package language

import "os"

// Language represents the sentences a language should have
type Language interface {
	Welcome() string
	Wrong() string
	Correct() string
	Challenge() string
	ID() string
}

// GetDefault returns the default selected language.
// Use the environment variable LANGUAGE to select a different language.
// Available languages: pt, en
func GetDefault() Language {
	langs := map[string]Language{
		"pt": pt{},
		"en": en{},
	}
	if l, ok := langs[os.Getenv("LANGUAGE")]; ok {
		return l
	}
	return langs["pt"]
}

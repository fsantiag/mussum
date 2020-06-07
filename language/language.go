package language

import "os"

type Language interface {
	Welcome() string
	Wrong() string
	Correct() string
	Challenge() string
}

func GetDefault() Language {
	langs := map[string]Language{
		"pt": Pt{},
		"en": En{},
	}
	if l, ok := langs[os.Getenv("LANGUAGE")]; ok {
		return l
	}
	return langs["pt"]
}

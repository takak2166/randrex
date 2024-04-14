package cmd

import (
	"regexp"
	"regexp/syntax"
	"testing"
)

func TestGenRandomStringFromRegex(t *testing.T) {
	patterns := []string{
		`^[a-z]{3}\d{2}[A-Z]{2}$`,
		`^hoge_.*`,
		`^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`,
		`https?://[\w/:%#\$&\?\(\)~\.=\+\-]+`,
	}

	for _, pattern := range patterns {
		re, err := syntax.Parse(pattern, syntax.Perl)
		if err != nil {
			t.Errorf("failed to parse pattern: %v", err)
			return
		}
		generatedStr := genRandomStringFromRegex(re)

		matched, err := regexp.MatchString(pattern, generatedStr)
		if err != nil {
			t.Errorf("error matching string: %v", err)
			return
		}

		if !matched {
			t.Errorf("generated string %q does not match expected pattern %q", generatedStr, pattern)
		}
	}
}

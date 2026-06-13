package text

import (
	"strings"
)

func Split(text, divider string) []string {
	if text == "" {
		return nil
	}
	return strings.Split(text, divider)
}

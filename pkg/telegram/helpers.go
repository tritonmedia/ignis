package telegram

import "strings"

// Proceed returns true if a text seems to be a yes
func Proceed(text string) bool {
	text = strings.ToLower(text)
	yesResults := locale.Proceed.Yes
	noResults := locale.Proceed.No

	yes := 0
	no := 0
	for _, s := range yesResults {
		if strings.Contains(text, s) {
			yes++
		}
	}

	for _, s := range noResults {
		if strings.Contains(text, s) {
			no++
		}
	}

	if yes > no {
		return true
	}

	return false
}

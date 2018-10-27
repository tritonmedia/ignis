package analysis

import (
	"log"
	"strings"
)

// ProceedAnalysis checks if we can proceed based on the results of the word used.
func ProceedAnalysis(msg string) bool {
	if model["proceed"][msg] == true {
		return true
	}

	tokens, err := tokenize(msg)
	if err != nil {
		log.Printf("[analysis/proceedanalysis:parse] failed to process: %s", err.Error())
		return false
	}

	scores := make([]bool, len(tokens))
	for i, token := range tokens {
		w := strings.ToLower(token.Text)
		scores[i] = model["proceed"][w]
	}

	pos := 0
	neg := 0
	for _, s := range scores {
		if s {
			pos++
		} else {
			neg++
		}
	}

	log.Printf("[analysis/proceedanalysis:score] pos=%d,neg=%d", pos, neg)

	if pos > neg {
		return true
	}

	return false
}

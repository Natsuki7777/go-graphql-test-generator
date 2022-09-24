package pkg

import (
	"regexp"
	"strings"
)

const (
	BRACE_OPEN   = "{"
	BRACE_CLOSE  = "}"
	PAREN_OPEN   = "("
	PAREN_CLOSE  = ")"
	SQUARE_OPEN  = "["
	SQUARE_CLOSE = "]"
	COLON        = ":"
	COMMA        = ","
	NEXT_LINE    = "\n"
)

func TokenizeInput(input string) []string {
	comment := regexp.MustCompile(`#.*`)
	input = comment.ReplaceAllString(input, "")
	input = strings.Replace(input, "!", "", -1)
	tokens := []string{}
	token := ""
	for _, c := range input {
		if string(c) == BRACE_OPEN || string(c) == BRACE_CLOSE || string(c) == PAREN_OPEN || string(c) == PAREN_CLOSE || string(c) == SQUARE_OPEN || string(c) == SQUARE_CLOSE || string(c) == COLON || string(c) == COMMA || string(c) == NEXT_LINE {
			if string(c) == NEXT_LINE {
				continue
			} else if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
			tokens = append(tokens, string(c))
		} else if string(c) == " " {
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			} else {
				continue
			}
		} else {
			token += string(c)
		}
	}
	return tokens
}

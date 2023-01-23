package tokenize

import "strings"

type EnTokenizer struct {
}

//Tokens splits text into tokens.
func (s *EnTokenizer) Tokens(text string) []string {
	return strings.Fields(text)
}

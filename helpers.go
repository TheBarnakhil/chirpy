package main

import (
	"slices"
	"strings"
)

// This function filters against the words "kerfuffle", "sharbert", "fornax" and returns "****" in their place
// It one arg the string on which the replace happens
// This function can be made more generic if we pass in the array of words to be replaced but that's going to be overkill for current scope
// The delimiter is set to a whitespace (" ")
// It returns the formatted string and also an error if any
func filterProfaneWords(og_str string) (string, error) {
	delimiter := " "

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(og_str, delimiter)
	newStr := []string{}
	for _, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			newStr = append(newStr, "****")
			continue
		}
		newStr = append(newStr, word)
	}
	return strings.Join(newStr, " "), nil
}

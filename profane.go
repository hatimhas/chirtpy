package main

import "strings"

// strings.ToLower
// strings.Split
// strings.Join
//
// words to filter
// kerfuffle
// sharbert
//
//	fornax
func profaneCheck(body string) string {
	bannedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	splitBody := strings.Split(body, " ")
	var stringBody []string
	for _, word := range splitBody {
		if _, found := bannedWords[strings.ToLower(word)]; found {
			stringBody = append(stringBody, "****")
		} else {
			stringBody = append(stringBody, word)
		}
	}
	cleanedBody := strings.Join(stringBody, " ")
	return cleanedBody
}

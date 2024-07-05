package tests

import "unicode"

// I exposed this function in order to be able to run it from test in separate test package
func GetDeepestLetter(input string) rune {
	var (
		deepestRune       rune = '?'
		currentLevel      int  = 0
		maxSavedRuneLevel int  = 0
	)
	for _, c := range input {
		switch {
		case c == '(':
			currentLevel++ // increase level
		case c == ')':
			currentLevel-- // decrease level
			// check if level is negative as that is clearly a malformed string
			if currentLevel < 0 {
				return '?'
			}
		case unicode.IsLower(c):
			// if current level is deeper than maxSavedRuneLevel
			if currentLevel > maxSavedRuneLevel {
				deepestRune = c
				maxSavedRuneLevel = currentLevel
			}
		}
	}
	if currentLevel != 0 {
		return '?'
	}
	return deepestRune
}

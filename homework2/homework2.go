package homework2

import (
	"strings"
	"unicode"
)

var vowels = "aeiou" // order matters

func ToPigLatin(s string) string {
	if s == "" {
		return s
	} else if isVowel(s[0]) {
		return s + "yay"
	} else {
		i := 0
		for isConsonant(s[i]) && i < len(s) {
			i++
		}
		return s[i:] + s[:i] + "ay"
	}
}

func isConsonant(character byte) bool {
	return !isVowel(character)
}

func isVowel(character byte) bool {
	return strings.ContainsRune(vowels, unicode.ToLower(rune(character)))
}

func EncodeVowels(s string) string {
	return strings.Map(func(character rune) rune {
		if index := strings.IndexRune(vowels, unicode.ToLower(character)); index != -1 {
			return rune('0' + index + 1)
		} else {
			return character
		}
	}, s)
}

func DecodeVowels(s string) string {
	return strings.Map(func(character rune) rune {
		if character >= '1' && character <= '5' {
			index := character - '0' - 1
			return rune(vowels[index])
		} else {
			return character
		}
	}, s)
}

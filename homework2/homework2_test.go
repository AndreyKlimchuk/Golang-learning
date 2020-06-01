package homework2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPigLatin(t *testing.T) {
	pairs := [][]string{
		{"pig", "igpay"},
		{"latin", "atinlay"},
		{"banana", "ananabay"},
		{"will", "illway"},
		{"butler", "utlerbay"},
		{"happy", "appyhay"},
		{"duck", "uckday"},
		{"me", "emay"},
		{"smile", "ilesmay"},
		{"string", "ingstray"},
		{"stupid", "upidstay"},
		{"glove", "oveglay"},
		{"trash", "ashtray"},
		{"floor", "oorflay"},
		{"store", "orestay"},
		{"eat", "eatyay"},
		{"omelet", "omeletyay"},
		{"are", "areyay"},
		{"egg", "eggyay"},
		{"explain", "explainyay"},
		{"always", "alwaysyay"},
		{"ends", "endsyay"},
		{"I", "Iyay"},
		{"", ""},
	}
	assertEqualPairs(t, pairs, ToPigLatin, "Conversion to Pig Latin should be correct.")
}

func TestVowelsEncoding(t *testing.T) {
	pairs := [][]string{
		{"pig", "p3g"},
		{"latin", "l1t3n"},
		{"Eat", "21t"},
		{"I", "3"},
	}
	assertEqualPairs(t, pairs, EncodeVowels, "Vowels encoding should be correct.")
}

func TestVowelsDecoding(t *testing.T) {
	pairs := [][]string{
		{"p3g", "pig"},
		{"l1t3n", "latin"},
		{"21t", "eat"},
		{"3", "i"},
	}
	assertEqualPairs(t, pairs, DecodeVowels, "Vowels decoding should be correct.")
}

func assertEqualPairs(t *testing.T, pairs [][]string, convert func(string) string, description string) {
	for _, v := range pairs {
		assert.Equal(t, v[1], convert(v[0]), description)
	}
}

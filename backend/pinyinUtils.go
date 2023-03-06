package backend

import (
	"fmt"
	"strings"
)

var pinyinTable map[string][]string = map[string][]string{
	"a": {"ā", "á", "ǎ", "à", "a"},
	"e": {"ē", "é", "ě", "è", "e"},
	"i": {"ī", "í", "ǐ", "ì", "i"},
	"o": {"ō", "ó", "ǒ", "ò", "o"},
	"u": {"ū", "ú", "ǔ", "ù", "u"},
	"ü": {"ǖ", "ǘ", "ǚ", "ǜ", "ü"},
}

func numberToTone(letter string, tone int) (string, error) {
	// cant fail since the letter passed in are all fixed
	tones, _ := pinyinTable[letter]
	return tones[tone-1], nil
}

func convertToTone(word string) (string, error) {
	lastByte := word[len(word)-1]
	// Simple ascii to int conversion
	tone := int(lastByte) - 48
	// Only works for bytes with 1,2,3,4,5
	if tone < 1 || tone > 5 {
		return word, nil
	}
	// We have a number so slice it off
	word = word[0 : len(word)-1]

	// From wikipedia:
	// If there is an a or an e, it will take the tone mark
	// If there is an ou, then the o takes the tone mark
	// Otherwise, the second vowel takes the tone mark

	// All possible vowel clusters
	// ài ào èi ià iào iè iò iù òu uà uài uè uì uò (üà) üè
	// Remove those with a or e
	// iò iù òu uì uò
	// Just need special cases for these 5 cases
	for _, test := range []struct {
		search  string
		replace string
	}{
		{"a", "a"},
		{"e", "e"},
		{"ou", "o"},
		{"io", "o"},
		{"iu", "u"},
		{"ui", "i"},
		{"uo", "o"},
		{"i", "i"},
		{"o", "o"},
		{"u", "u"},
		{"ü", "ü"},
	} {
		if strings.Contains(word, test.search) {
			replacement, err := numberToTone(test.replace, tone)
			if err != nil {
				return "", err
			}
			return strings.Replace(word, test.replace, replacement, -1), nil
		}
	}

	return "", fmt.Errorf("No tone match")
}

// For now we just worry about the format we expect from cc-cedict
func ToStandardPinyin(pinyin string) string {
	if pinyin == "" {
		return pinyin
	}
	pinyin = strings.Replace(pinyin, "u:", "ü", -1)
	pinyinSlice := strings.Split(pinyin, " ")

	builder := strings.Builder{}
	for _, word := range pinyinSlice {
		converted, err := convertToTone(word)
		if err != nil {
			builder.WriteString(word)
		} else {
			// For now failures just rewrite the original
			builder.WriteString(converted)
		}
	}
	return builder.String()
}

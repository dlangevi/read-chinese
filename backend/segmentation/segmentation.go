package segmentation

import (
	"os"
	"regexp"
	"strings"

	"github.com/yanyiwu/gojieba"
)

type Segmentation struct {
}

type Token struct {
	Data   string
	IsWord bool
}

var jieba *gojieba.Jieba
var punctuation *regexp.Regexp
var whitespace *regexp.Regexp
var latin *regexp.Regexp
var chinese *regexp.Regexp

func isChineseWord(word string) bool {
	if punctuation.MatchString(word) {
		return false
	}
	if whitespace.MatchString(word) {
		return false
	}
	if latin.MatchString(word) {
		return false
	}
	if chinese.MatchString(word) {
		return true
	}
	return false
}

// segmentSentence
func (s *Segmentation) SegmentSentence(sentence string) []Token {
	words := jieba.Cut(sentence, false)
	tokens := []Token{}
	for _, word := range words {
		isWord := isChineseWord(word)
		tokens = append(tokens, Token{
			Data:   word,
			IsWord: isWord,
		})
	}
	return tokens
}

type FrequencyTable map[string]int

// doFullSegmentation
func (s *Segmentation) SegmentFullText(path string) ([]string, FrequencyTable, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	words := jieba.Cut(string(bytes), false)
	sentences := []string{}
	frequency := FrequencyTable{}
	previousSentence := strings.Builder{}
	currentSentence := strings.Builder{}

	incrementFrequency := func(word string) {
		current, ok := frequency[word]
		if !ok {
			current = 0
		}
		current = current + 1
		frequency[word] = current
	}

	terminateSentence := func() {
		if previousSentence.Len() > 0 {
			sentences = append(sentences, previousSentence.String())
			previousSentence.Reset()
		}
		previousSentence.WriteString(currentSentence.String())
		currentSentence.Reset()
	}

	for _, word := range words {
		isWord := isChineseWord(word)
		if isWord {
			incrementFrequency(word)
		}

		// 15. 14. etc can be tokens. These used to break storage?
		if len(word) > 1 && strings.Contains(word, ".") {
			word = strings.ReplaceAll(word, ".", "")
		}
		if word == "\n" {
			terminateSentence()
		} else if word == "？" || word == "！" || word == "。" || word == "…" || word == "." {
			if currentSentence.Len() == 0 {
				previousSentence.WriteString(word)
			} else {
				currentSentence.WriteString(word)
				terminateSentence()
			}
		} else if word == " " || word == "　" || word == "\t" {
			// Do I actually need this?
			if currentSentence.Len() > 0 {
				currentSentence.WriteString(word)
			}
		} else if word == "”" || word == "‘" || word == "』" {
			// Closing quotes go onto previous sentence if needed
			if currentSentence.Len() == 0 {
				previousSentence.WriteString(word)
			} else {
				currentSentence.WriteString(word)
			}
		} else {
			currentSentence.WriteString(word)
		}
	}

	if previousSentence.Len() > 0 {
		sentences = append(sentences, previousSentence.String())
	}
	if currentSentence.Len() > 0 {
		sentences = append(sentences, currentSentence.String())
	}

	return sentences, frequency, nil

}

// todo
// computeDict (modify the default dict to only use words from user dicts)
// loadJieba
func NewSegmentation() (*Segmentation, error) {
	s := &Segmentation{}
	jieba = gojieba.NewJieba()
	var err error
	punctuation, err = regexp.Compile(`\p{P}`)
	if err != nil {
		return s, err
	}
	whitespace, err = regexp.Compile(`\s+`)
	if err != nil {
		return s, err
	}
	latin, err = regexp.Compile(`\p{Latin}`)
	if err != nil {
		return s, err
	}
	chinese, err = regexp.Compile(`\p{Han}`)
	if err != nil {
		return s, err
	}
	return s, nil
}

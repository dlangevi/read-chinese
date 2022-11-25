package backend

import (
	"bufio"
	"embed"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/yanyiwu/gojieba"
)

type Segmentation struct {
}

//go:embed assets/dict
var jiebaDicts embed.FS
var jiebaFiles = map[string]string{
	"DICT_PATH":       "jieba.dict.utf8",
	"HMM_PATH":        "hmm_model.utf8",
	"USER_DICT_PATH":  "user.dict.utf8",
	"IDF_PATH":        "idf.utf8",
	"STOP_WORDS_PATH": "stop_words.utf8",
}

func UnloadJiebaDicts() error {
	for _, filename := range jiebaFiles {
		data, _ := jiebaDicts.ReadFile(path.Join("assets", "dict", filename))
		dest := ConfigDir("jiebaDicts", filename)
		err := os.WriteFile(dest, data, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

type Token struct {
	Data   string
	IsWord bool
}

var jieba *gojieba.Jieba
var punctuation = regexp.MustCompile(`\p{P}`)
var whitespace = regexp.MustCompile(`\s+`)
var latin = regexp.MustCompile(`\p{Latin}`)
var chinese = regexp.MustCompile(`\p{Han}`)

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

func constructDict(d *Dictionaries) error {
	dict, err := jiebaDicts.Open(path.Join(
		"assets", "dict", jiebaFiles["DICT_PATH"]))
	if err != nil {
		return err
	}
	defer dict.Close()
	replacementPath := path.Join(os.TempDir(), "replacement.dict.utf8")
	userDict, err := os.Create(replacementPath)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(dict)
	wr := bufio.NewWriter(userDict)
	totalWords := 0
	validWords := 0
	for sc.Scan() {
		line := sc.Text()
		parts := strings.Split(line, " ")
		word := parts[0]
		totalWords += 1
		if d.IsInDictionary(word) {
			validWords += 1
			wr.WriteString(line)
			wr.WriteRune('\n')
		}
	}
	wr.Flush()
	userDict.Close()
	// If the user has not installed dictionaries we just use the
	// default segmentation
	if validWords != 0 {
		jieba = gojieba.NewJieba(
			replacementPath,
			ConfigDir("jiebaDicts", jiebaFiles["HMM_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["USER_DICT_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["IDF_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["STOP_WORDS_PATH"]),
		)
	} else {
		jieba = gojieba.NewJieba(
			ConfigDir("jiebaDicts", jiebaFiles["DICT_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["HMM_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["USER_DICT_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["IDF_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["STOP_WORDS_PATH"]),
		)
	}
	log.Println("totalWords", totalWords, "validWords", validWords)
	return nil
}

func NewSegmentation(d *Dictionaries) (*Segmentation, error) {
	s := &Segmentation{}
	err := constructDict(d)
	if err != nil {
		return s, err
	}
	return s, nil
}

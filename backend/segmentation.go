package backend

import (
	"bufio"
	"embed"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yanyiwu/gojieba"
)

type Segmentation struct {
	jieba *gojieba.Jieba
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

type TokenizedSentence []Token

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
	words := s.jieba.Cut(sentence, false)
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
func (s *Segmentation) SegmentFullText(path string) ([]TokenizedSentence, FrequencyTable, error) {

	// TODO while we have both types of paths
	if !filepath.IsAbs(path) {
		path = ConfigDir(path)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	words := s.jieba.Cut(string(bytes), false)
	segmentedSentences := []TokenizedSentence{}
	previousSegmented := TokenizedSentence{}
	currentSegmented := TokenizedSentence{}
	frequency := FrequencyTable{}

	incrementFrequency := func(word string) {
		current, ok := frequency[word]
		if !ok {
			current = 0
		}
		current = current + 1
		frequency[word] = current
	}

	terminateSentence := func() {
		if len(previousSegmented) > 0 {
			segmentedSentences = append(segmentedSentences, previousSegmented)
			previousSegmented = TokenizedSentence{}
		}
		previousSegmented = currentSegmented
		currentSegmented = TokenizedSentence{}
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
			if len(currentSegmented) == 0 {
				previousSegmented = append(previousSegmented, Token{
					Data:   word,
					IsWord: isWord,
				})
			} else {
				currentSegmented = append(currentSegmented, Token{
					Data:   word,
					IsWord: isWord,
				})
				terminateSentence()
			}
		} else if word == " " || word == "　" || word == "\t" {
			// Do I actually need this?
			if len(currentSegmented) > 0 {
				currentSegmented = append(currentSegmented, Token{
					Data:   word,
					IsWord: isWord,
				})
			}
		} else if word == "”" || word == "‘" || word == "』" {
			// Closing quotes go onto previous sentence if needed
			if len(currentSegmented) == 0 {
				previousSegmented = append(previousSegmented, Token{
					Data:   word,
					IsWord: isWord,
				})
			} else {
				currentSegmented = append(currentSegmented, Token{
					Data:   word,
					IsWord: isWord,
				})
			}
		} else {
			currentSegmented = append(currentSegmented, Token{
				Data:   word,
				IsWord: isWord,
			})
		}
	}

	if len(previousSegmented) > 0 {
		segmentedSentences = append(segmentedSentences, previousSegmented)
	}
	if len(currentSegmented) > 0 {
		segmentedSentences = append(segmentedSentences, currentSegmented)
	}

	return segmentedSentences, frequency, nil
}

func toString(sentence []Token) string {
	wordSentence := strings.Builder{}
	for _, token := range sentence {
		wordSentence.WriteString(token.Data)
	}
	return wordSentence.String()
}

func constructDict(d *Dictionaries, s *Segmentation) error {
	dict, err := jiebaDicts.Open(path.Join(
		"assets", "dict", jiebaFiles["DICT_PATH"]))
	if err != nil {
		return err
	}
	defer dict.Close()
	replacementPath := ConfigDir("jiebaDicts", "replacement.dict.utf8")
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
		s.jieba = gojieba.NewJieba(
			replacementPath,
			ConfigDir("jiebaDicts", jiebaFiles["HMM_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["USER_DICT_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["IDF_PATH"]),
			ConfigDir("jiebaDicts", jiebaFiles["STOP_WORDS_PATH"]),
		)
	} else {
		s.jieba = gojieba.NewJieba(
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

func (s *Segmentation) ReloadJieba(d *Dictionaries) error {
	err := constructDict(d, s)
	if err != nil {
		return err
	}
	return nil
}

func NewSegmentation(d *Dictionaries) (*Segmentation, error) {
	s := &Segmentation{}
	err := constructDict(d, s)
	if err != nil {
		return s, err
	}
	d.PassSegmentation(s)
	return s, nil
}

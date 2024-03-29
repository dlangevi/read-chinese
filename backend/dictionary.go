package backend

import (
	"compress/gzip"
	"encoding/json"
	// "log"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const URL = "https://www.mdbg.net/chinese/export/cedict/cedict_1_0_ts_utf-8_mdbg.txt.gz"

type (
	Dictionary interface {
		GetDefinitions(word string) []string
		GetEntries(word string) []DictionaryDefinition
		GetPronuciations(word string) []string
		GetPartialMatches(partial string) []string
		Contains(word string) bool
	}

	dictionary struct {
		Definitions map[string][]dictionaryEntry
	}

	dictionaryEntry struct {
		Definition    string
		Pronunciation string
	}

	// The type that interfaces with typescript front end
	DictionaryDefinition struct {
		Definition    string `json:"definition"`
		Pronunciation string `json:"pronunciation,omitempty"`
	}
)

func newDictionary() *dictionary {
	return &dictionary{
		Definitions: map[string][]dictionaryEntry{},
	}
}

func SaveDictionary(d *dictionary, path string) error {
	rawDictBytes, err := json.Marshal(d)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, rawDictBytes, 0666)
	return err
}

var entryRegexp = regexp.MustCompile(`^(.*) (.*) \[(.*)\] /(?:(.*)/)+`)

func downloadCedict(backend *Backend) (io.ReadCloser, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Printf("Predicting %v length\n", resp.ContentLength)

	counter := &ReportedDownload{
		backend:    backend,
		TotalBytes: uint64(resp.ContentLength),
	}
	gz, err := gzip.NewReader(io.TeeReader(resp.Body, counter))
	if err != nil {
		return nil, err
	}

	return gz, nil
}

func FromCedictFormat(backend *Backend) (*dictionary, error) {
	dictionary := newDictionary()
	ccdict, err := downloadCedict(backend)
	if err != nil {
		return dictionary, err
	}
	sc := bufio.NewScanner(ccdict)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") { // Ignore comments
			continue
		}
		matches := entryRegexp.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("Entry doesn't match regular expression")
		}
		// e.Traditional = matches[1]
		simplified := matches[2]
		pinyin := strings.ToLower(matches[3])
		defs := matches[4]
		// defs := strings.Split(matches[4], "/")
		// if len(defs) == 0 {
		// 	return nil, fmt.Errorf("No definitions found")
		// }
		terms, ok := dictionary.Definitions[simplified]
		if !ok {
			terms = []dictionaryEntry{}
		}
		terms = append(terms, dictionaryEntry{
			Definition:    defs,
			Pronunciation: pinyin,
		})
		dictionary.Definitions[simplified] = terms
	}

	return dictionary, nil
}

func FromSavedDictionary(path string) (*dictionary, error) {
	dictionary := newDictionary()
	rawDictBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawDictBytes, dictionary)
	if err != nil {
		return nil, err
	}

	return dictionary, nil
}

func FromMigakuJsonFormat(path string) (*dictionary, error) {
	type MigakuDictionaryEntry struct {
		Term          string
		Pronunciation string
		Definition    string
	}
	type MigakuDictionary struct {
		Entries []MigakuDictionaryEntry
	}
	dictionary := newDictionary()

	rawDict := &MigakuDictionary{}
	rawDictBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawDictBytes, &rawDict.Entries)
	if err != nil {
		return nil, err
	}

	for _, entry := range rawDict.Entries {
		terms, ok := dictionary.Definitions[entry.Term]
		if !ok {
			terms = []dictionaryEntry{}
		}
		terms = append(terms, dictionaryEntry{
			Definition:    entry.Definition,
			Pronunciation: entry.Pronunciation,
		})
		dictionary.Definitions[entry.Term] = terms
	}

	return dictionary, nil
}

func (d *dictionary) GetDefinitions(word string) []string {
	definitions := []string{}
	entries, ok := d.Definitions[word]
	if !ok {
		return definitions
	}
	for _, entry := range entries {
		definitions = append(definitions, entry.Definition)
	}
	return definitions
}

func (d *dictionary) GetEntries(word string) []DictionaryDefinition {
	definitions := []DictionaryDefinition{}

	entries, ok := d.Definitions[word]
	if !ok {
		return definitions
	}

	for _, entry := range entries {
		definitions = append(definitions, DictionaryDefinition{
			Definition:    entry.Definition,
			Pronunciation: ToStandardPinyin(entry.Pronunciation),
		})
	}

	return definitions
}

func (d *dictionary) GetPronuciations(word string) []string {
	pronuciation := []string{}
	entries, ok := d.Definitions[word]
	if !ok {
		return pronuciation
	}
	for _, entry := range entries {
		pronuciation = append(pronuciation, ToStandardPinyin(entry.Pronunciation))
	}
	return pronuciation
}

func (d *dictionary) GetPartialMatches(partial string) []string {
	words := []string{}
	for word := range d.Definitions {
		if strings.Contains(word, partial) {
			words = append(words, word)
		}
	}
	return words
}

func (d *dictionary) Contains(word string) bool {
	_, ok := d.Definitions[word]
	return ok
}

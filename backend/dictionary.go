package backend

import (
	"encoding/json"
	// "log"
	"os"
	"strings"
)

type (
	Dictionary interface {
		GetDefinitions(word string) []string
		GetEntries(word string) []DictionaryDefinition
		GetPronuciations(word string) []string
		GetPartialMatches(partial string) []string
		Contains(word string) bool
	}

	dictionary struct {
		definitions map[string][]dictionaryEntry
	}

	dictionaryEntry struct {
		Definition    string
		Pronunciation string
	}

	// The type that interfaces with typescript front end
	DictionaryDefinition struct {
		Definition    string `json:"definition"`
		Pronunciation string `json:"pronunciation"`
	}
)

func newDictionary() *dictionary {
	return &dictionary{
		definitions: map[string][]dictionaryEntry{},
	}
}

func FromMikaguJsonFormat(path string) (*dictionary, error) {
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
		terms, ok := dictionary.definitions[entry.Term]
		if !ok {
			terms = []dictionaryEntry{}
		}
		terms = append(terms, dictionaryEntry{
			Definition:    entry.Definition,
			Pronunciation: entry.Pronunciation,
		})
		dictionary.definitions[entry.Term] = terms
	}

	return dictionary, nil
}

func (d *dictionary) GetDefinitions(word string) []string {
	definitions := []string{}
	entries, ok := d.definitions[word]
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

	entries, ok := d.definitions[word]
	if !ok {
		return definitions
	}

	for _, entry := range entries {
		definitions = append(definitions, DictionaryDefinition{
			Definition:    entry.Definition,
			Pronunciation: entry.Pronunciation,
		})
	}

	return definitions
}

func (d *dictionary) GetPronuciations(word string) []string {
	pronuciation := []string{}
	entries, ok := d.definitions[word]
	if !ok {
		return pronuciation
	}
	for _, entry := range entries {
		pronuciation = append(pronuciation, entry.Pronunciation)
	}
	return pronuciation
}

func (d *dictionary) GetPartialMatches(partial string) []string {
	words := []string{}
	for word := range d.definitions {
		if strings.Contains(word, partial) {
			words = append(words, word)
		}
	}
	return words
}

func (d *dictionary) Contains(word string) bool {
	_, ok := d.definitions[word]
	return ok
}

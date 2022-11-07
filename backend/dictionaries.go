package backend

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type DictionaryEntry struct {
	Definition    string `json:"definition"`
	Pronunciation string `json:"pronunciation"`
}

type Dictionary struct {
	Definitions map[string][]DictionaryEntry
	Language    string
}

type Dictionaries struct {
	PrimaryDictName string
	PrimaryDict     *Dictionary
	Dictionaries    map[string]*Dictionary
}

type RawDictionaryEntry struct {
	Term          string
	Pronunciation string
	Definition    string
}

type RawDictionary struct {
	Entries []RawDictionaryEntry
}

func NewDictionaries() *Dictionaries {
	dicts := &Dictionaries{}
	dicts.loadDictionaries()
	primaryName := userSettings.PrimaryDict
	dicts.PrimaryDictName = primaryName
	// TODO this could fail
	dicts.PrimaryDict, _ = dicts.Dictionaries[primaryName]

	return dicts
}

func parseDictionaryFile(path string) (*RawDictionary, error) {
	rawDict := &RawDictionary{}
	rawDictBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(rawDictBytes, &rawDict.Entries)

	return rawDict, nil
}

func (d *Dictionaries) loadDictionaries() error {
	d.Dictionaries = map[string]*Dictionary{}
	return d.transformRawDictionaries(userSettings.Dicts)
}

func (d *Dictionaries) transformRawDictionaries(dicts map[string]Dict) error {
	for name, dict := range dicts {
		newDict := Dictionary{
			Definitions: map[string][]DictionaryEntry{},
			Language:    dict.Language,
		}
		rawDict, err := parseDictionaryFile(dict.Path)
		if err != nil {
			return err
		}

		for _, entry := range rawDict.Entries {
			terms, ok := newDict.Definitions[entry.Term]
			if !ok {
				terms = []DictionaryEntry{}
			}
			terms = append(terms, DictionaryEntry{
				Definition:    entry.Definition,
				Pronunciation: entry.Pronunciation,
			})
			newDict.Definitions[entry.Term] = terms
		}
		d.Dictionaries[name] = &newDict

	}
	return nil
}

func (d *Dictionaries) AddDictionary(name string, path string, language string) {
	SaveDict(name, path, language)
	d.loadDictionaries()
}

func (d *Dictionaries) DeleteDictionary(name string) {
	DeleteDict(name)
	delete(d.Dictionaries, name)
}

func (d *Dictionaries) SetPrimaryDict(name string) {
	d.PrimaryDictName = name
	// TODO this could fail
	primary, ok := d.Dictionaries[name]
	if !ok {
		log.Fatal("not ok dict")

	}

	d.PrimaryDict = primary
	SetPrimaryDict(name)
}

type DictionaryInfoMap map[string]DictionaryInfo
type DictionaryInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Language string `json:"type"`
}

func (d *Dictionaries) GetDictionaryInfo() DictionaryInfoMap {
	dictInfoMap := DictionaryInfoMap{}
	for name, dict := range userSettings.Dicts {
		dictInfoMap[name] = DictionaryInfo{
			Name:     name,
			Path:     dict.Path,
			Language: dict.Language,
		}
	}
	return dictInfoMap
}

func (d *Dictionaries) GetDefinitionsForWord(word string, language string) []DictionaryEntry {
	answers := []DictionaryEntry{}

	for _, dict := range d.Dictionaries {
		if dict.Language != language {
			continue
		}
		terms, ok := dict.Definitions[word]
		if !ok {
			continue
		}
		answers = append(answers, terms...)
	}
	return answers
}

type UnknownWordEntry struct {
	Word       string `json:"word"`
	Occurance  int    `json:"occurance,omitempty"`
	Definition string `json:"definition,omitempty"`
	Pinyin     string `json:"pinyin,omitempty"`
}

func (d *Dictionaries) getDefaultDefinition(word string) string {
	term, ok := d.PrimaryDict.Definitions[word]
	if !ok {
		return ""
	}
	first := term[0]
	return first.Definition
}

func (d *Dictionaries) getPinyin(word string) string {
	term, ok := d.PrimaryDict.Definitions[word]
	if !ok {
		return ""
	}
	// TODO look into go-funk if more occasion for this stuff arrises
	pinyin := []string{}
	for _, item := range term {
		pinyin = append(pinyin, item.Pronunciation)
	}

	return strings.Join(pinyin, ", ")
}

// TODO this passing back and forth of UnknownWordEntry feels clunky
func (d *Dictionaries) GetDefinitions(words []UnknownWordEntry) []UnknownWordEntry {
	for index := range words {
		word := &words[index]
		word.Definition = d.getDefaultDefinition(word.Word)
		word.Pinyin = d.getPinyin(word.Word)
	}
	return words
}

func (d *Dictionaries) GetPossibleWords(partial string) []UnknownWordEntry {
	words := []UnknownWordEntry{}

	for word := range d.PrimaryDict.Definitions {
		if strings.Contains(word, partial) && !known.isKnown(word) {
			words = append(words, UnknownWordEntry{
				Word: word,
			})
		}
	}

	return d.GetDefinitions(words)
}

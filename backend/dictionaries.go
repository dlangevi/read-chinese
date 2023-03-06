package backend

import (
	"log"
	"strings"
)

type Dictionaries struct {
	PrimaryDictName string
	PrimaryDict     Dictionary
	Dictionaries    map[string]UserDictionary
	userSettings    *UserConfig
	known           *KnownWords
	Segmentation    *Segmentation
	backend         *Backend
}

type UserDictionary struct {
	Dictionary Dictionary
	Language   string
}

func NewDictionaries(
	backend *Backend,
	userSettings *UserConfig,
	known *KnownWords,
) *Dictionaries {
	dicts := &Dictionaries{
		backend:      backend,
		userSettings: userSettings,
		known:        known,
	}
	dicts.loadDictionaries()
	primaryName := userSettings.DictionariesConfig.PrimaryDict
	dicts.PrimaryDictName = primaryName
	// TODO this could fail
	dicts.PrimaryDict = dicts.Dictionaries[primaryName].Dictionary

	return dicts
}

func (d *Dictionaries) HealthCheck() string {
	if len(d.Dictionaries) == 0 {
		return "User has no dictionaries"
	}
	return ""
}

// Is this the way to do a reference loop?
func (d *Dictionaries) PassSegmentation(s *Segmentation) {
	d.Segmentation = s
}

func (d *Dictionaries) loadDictionaries() error {
	d.Dictionaries = map[string]UserDictionary{}
	for name, dict := range d.userSettings.DictionariesConfig.Dicts {
		newDict, err := FromSavedDictionary(dict.Path)
		if err != nil {
			return err
		}
		d.Dictionaries[name] = UserDictionary{
			newDict,
			dict.Language,
		}
	}
	primaryName := d.userSettings.DictionariesConfig.PrimaryDict
	d.PrimaryDictName = primaryName
	primaryDict, ok := d.Dictionaries[primaryName]
	if !ok {
		if len(d.Dictionaries) > 0 {
			for name := range d.Dictionaries {
				primaryName = name
				d.SetPrimaryDict(name)
				continue
			}
		}
	} else {
		d.PrimaryDict = primaryDict.Dictionary
	}
	if d.Segmentation != nil {
		err := d.Segmentation.ReloadJieba(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dictionaries) AddCedict() {
	// Parse it into a *dictionary
	d.backend.setupProgress("Downloading cc-cedict", 100)
	dictionary, err := FromCedictFormat(d.backend)
	if err != nil {
		log.Println(err)
	}
	savedPath := ConfigDir("userDicts", "cc-cedict")
	SaveDictionary(dictionary, savedPath)
	d.userSettings.SaveDict("cc-cedict", savedPath, "english")
	d.loadDictionaries()
}

func (d *Dictionaries) AddMigakuDictionary(name string, path string, language string) {
	// TODO error handle
	dictionary, _ := FromMigakuJsonFormat(path)
	savedPath := ConfigDir("userDicts", name)
	SaveDictionary(dictionary, savedPath)
	d.userSettings.SaveDict(name, savedPath, language)
	d.loadDictionaries()
}

func (d *Dictionaries) DeleteDictionary(name string) {
	d.userSettings.DeleteDict(name)
	delete(d.Dictionaries, name)
}

func (d *Dictionaries) SetPrimaryDict(name string) {
	d.PrimaryDictName = name
	// TODO this could fail
	primary, ok := d.Dictionaries[name]
	if !ok {
		log.Fatal("not ok dict")

	}

	d.PrimaryDict = primary.Dictionary
	d.userSettings.SetPrimaryDict(name)
}

type DictionaryInfoMap map[string]DictionaryInfo
type DictionaryInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Language  string `json:"type"`
	IsPrimary bool   `json:"isPrimary"`
}

func (d *Dictionaries) ExportDictionaryInfo() DictionaryInfo {
	return DictionaryInfo{}
}

func (d *Dictionaries) GetDictionaryInfo() DictionaryInfoMap {
	dictInfoMap := DictionaryInfoMap{}
	for name, dict := range d.userSettings.DictionariesConfig.Dicts {
		dictInfoMap[name] = DictionaryInfo{
			Name:      name,
			Path:      dict.Path,
			Language:  dict.Language,
			IsPrimary: name == d.PrimaryDictName,
		}
	}
	return dictInfoMap
}

func (d *Dictionaries) GetDefinitionsForWord(word string, language string) []DictionaryDefinition {
	answers := []DictionaryDefinition{}

	for _, dict := range d.Dictionaries {
		if dict.Language != language {
			continue
		}
		terms := dict.Dictionary.GetEntries(word)
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
	terms := d.PrimaryDict.GetDefinitions(word)
	if len(terms) == 0 {
		return ""
	}
	return terms[0]
}

func (d *Dictionaries) getPinyin(word string) string {
	pronuciationMap := map[string]struct{}{}
	for _, dict := range d.Dictionaries {
		proniciations := dict.Dictionary.GetPronuciations(word)
		for _, pronuciation := range proniciations {
			pronuciationMap[pronuciation] = struct{}{}
		}
	}

	pinyin := []string{}
	for pronuciation := range pronuciationMap {
		pinyin = append(pinyin, pronuciation)
	}

	return strings.Join(pinyin, ", ")
}

// TODO this passing back and forth of UnknownWordEntry feels clunky
type WordDefinitions map[string]DictionaryDefinition

func (d *Dictionaries) GetDefinitions(words []string) WordDefinitions {
	entries := WordDefinitions{}
	for _, word := range words {
		entries[word] = DictionaryDefinition{
			Definition:    d.getDefaultDefinition(word),
			Pronunciation: d.getPinyin(word),
		}
	}
	return entries
}

func (d *Dictionaries) GetPossibleWords(partial string) []UnknownWordEntry {
	words := d.PrimaryDict.GetPartialMatches(partial)
	unknown := []UnknownWordEntry{}
	// TODO filter out known
	for _, word := range words {
		unknown = append(unknown, UnknownWordEntry{
			Word: word,
		})
	}
	return unknown
}

func (d *Dictionaries) IsInDictionary(word string) bool {
	for _, dict := range d.Dictionaries {
		ok := dict.Dictionary.Contains(word)
		if ok {
			return true
		}
	}
	return false
}

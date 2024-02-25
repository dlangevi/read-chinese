package backend

import (
	"errors"
	"log"
	"path"
	"strings"
)

type Dictionaries struct {
	PrimaryDictName string
	PrimaryDict     Dictionary
	Dictionaries    map[string]UserDictionary
	backend         *Backend
}

type UserDictionary struct {
	Dictionary Dictionary
	Language   string
}

func NewDictionaries(
	backend *Backend,
) *Dictionaries {
	dicts := &Dictionaries{
		backend: backend,
	}
	dicts.loadDictionaries()
	primaryName := backend.UserSettings.DictionariesConfig.PrimaryDict
	dicts.PrimaryDictName = primaryName
	// TODO this could fail
	dicts.PrimaryDict = dicts.Dictionaries[primaryName].Dictionary

	return dicts
}

func (d *Dictionaries) HealthCheck() error {
	if len(d.Dictionaries) == 0 {
		return errors.New("User has no dictionaries")
	}
	return nil
}

// Is this the way to do a reference loop?
func (d *Dictionaries) PassSegmentation(s *Segmentation) {
	d.backend.Segmentation = s
}

func (d *Dictionaries) loadDictionaries() error {
	d.Dictionaries = map[string]UserDictionary{}
	for name, dict := range d.backend.UserSettings.DictionariesConfig.Dicts {
		dictPath := dict.Path
		if !path.IsAbs(dict.Path) {
			dictPath = ConfigDir(dictPath)
		}
		newDict, err := FromSavedDictionary(dictPath)
		if err != nil {
			return err
		}
		d.Dictionaries[name] = UserDictionary{
			newDict,
			dict.Language,
		}
	}
	primaryName := d.backend.UserSettings.DictionariesConfig.PrimaryDict
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
	if d.backend.Segmentation != nil {
		err := d.backend.Segmentation.ReloadJieba(d)
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
	relPath := path.Join("userDicts", "cc-cedict")
	savedPath := ConfigDir(relPath)
	SaveDictionary(dictionary, savedPath)
	d.backend.UserSettings.SaveDict("cc-cedict", relPath, "english")
	d.loadDictionaries()
}

func (d *Dictionaries) AddMigakuDictionary(name string, dictpath string, language string) {
	// TODO error handle
	dictionary, _ := FromMigakuJsonFormat(dictpath)
	relPath := path.Join("userDicts", name)
	savedPath := ConfigDir(relPath)
	SaveDictionary(dictionary, savedPath)
	d.backend.UserSettings.SaveDict(name, relPath, language)
	d.loadDictionaries()
}

func (d *Dictionaries) DeleteDictionary(name string) {
	d.backend.UserSettings.DeleteDict(name)
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
	d.backend.UserSettings.SetPrimaryDict(name)
}

type DictionaryInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Language  string `json:"type"`
	IsPrimary bool   `json:"isPrimary"`
}

func (d *Dictionaries) ExportDictionaryInfo() DictionaryInfo {
	return DictionaryInfo{}
}

func (d *Dictionaries) GetDictionaryInfo() map[string]DictionaryInfo {
	dictInfoMap := map[string]DictionaryInfo{}
	for name, dict := range d.backend.UserSettings.DictionariesConfig.Dicts {
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

// Wails auto typescript generation fails on this so use WordDefinitions
// To avoid issues
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

func (d *Dictionaries) GetPossibleWords(partial string) []string {
	words := d.PrimaryDict.GetPartialMatches(partial)
	return words
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

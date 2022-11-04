package backend

import (
	"testing"
)

func TestParseDictionaryFile(t *testing.T) {
	dict, err := parseDictionaryFile("testdata/testdata.json")
	if err != nil {
		t.Errorf("Failed to load raw dictionary: %v", err)
	}
	if len(dict.Entries) != 7 {
		t.Errorf("Raw dict had wrong number of entries: %v", len(dict.Entries))
	}
}

func TestTransformRawDictionaries(t *testing.T) {
	dicts := &Dictionaries{
		Dictionaries: map[string]*Dictionary{},
	}
	dicts.transformRawDictionaries(map[string]Dict{
		"test": {
			Path:     "testdata/testdata.json",
			Language: "english",
		},
	})
	dict, ok := dicts.Dictionaries["test"]
	if !ok {
		t.Errorf("Dictionary test failed to be loaded")
	}
	if len(dict.Definitions) != 6 {
		t.Errorf("Dict had wrong number of entries: %v", len(dict.Definitions))
	}

	好, ok := dict.Definitions["好"]
	if len(好) != 2 {
		t.Errorf("Did not group definitions for 好")
	}
}

func loadTestDictionaries(d *Dictionaries) error {
	d.Dictionaries = map[string]*Dictionary{}
	return d.transformRawDictionaries(map[string]Dict{
		"test": {
			Path:     "testdata/testdata.json",
			Language: "english",
		},
	})
}

func TestDefinitions(t *testing.T) {
	dicts := &Dictionaries{}
	loadTestDictionaries(dicts)
	dicts.PrimaryDictName = "test"
	dictPtr, ok := dicts.Dictionaries["test"]
	if !ok {
		t.Errorf("Failed to initialize primary dict")
	}
	dicts.PrimaryDict = dictPtr
	defs := dicts.getDefaultDefinition("善良")
	if defs != "good and honest/kindhearted" {
		t.Errorf("Wrong definition %v", defs)
	}
	pinyin := dicts.getPinyin("好")
	if pinyin != "hǎo, hào" {
		t.Errorf("Wrong pinyin %v", pinyin)
	}

}

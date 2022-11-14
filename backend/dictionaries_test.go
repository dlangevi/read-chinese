package backend

import (
	"testing"
)

// Make sure the dictionary loaded from test data is doing the correct
// Data grouping
func TestDictionaryContents(t *testing.T) {

	dicts := testRuntime.Dictionaries

	dict, ok := dicts.Dictionaries["example"]
	if !ok {
		t.Errorf("Dictionary example failed to be loaded")
	}
	if len(dict.Definitions) != 83 {
		t.Errorf("Dict had wrong number of entries: %v", len(dict.Definitions))
	}

	好, ok := dict.Definitions["好"]
	if len(好) != 2 {
		t.Errorf("Did not group definitions for 好")
	}

	你好, ok := dict.Definitions["你好"]
	if !ok || len(你好) != 1 {
		t.Errorf("Did not group definitions for 你好")
	}

	地方, ok := dict.Definitions["地方"]
	if !ok || len(地方) != 2 {
		t.Errorf("Did not group definitions for 地方")
	}

	的, ok := dict.Definitions["的"]
	if !ok || len(的) != 4 {
		t.Errorf("Did not group definitions for 的")
	}
}

func TestGetDictionaryInfo(t *testing.T) {
	info := testRuntime.Dictionaries.GetDictionaryInfo()
	if len(info) != 1 {
		t.Errorf("Info doesn't provide dictionary")
	}
	dict, ok := info["example"]
	if !ok {
		t.Errorf("Example dictionary served with wrong name")
	}
	if dict.Language != "english" || dict.Name != "example" {
		t.Errorf("Example dictionary served with wrong data")
	}
}

func TestGetDefinitions(t *testing.T) {
	defs := testRuntime.Dictionaries.GetDefinitionsForWord("的", "english")
	if len(defs) != 4 {
		t.Errorf("Not enough defs")
	}

	possible := testRuntime.Dictionaries.GetPossibleWords("你")
	if len(possible) != 5 {
		t.Errorf("Not enough options %v", possible)
	}

}

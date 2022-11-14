package backend

import (
	"testing"
)

// Make sure the dictionary loaded from test data is doing the correct
// Data grouping
func TestDictionaryContents(t *testing.T) {

	dicts := runtime.Dictionaries

	dict, ok := dicts.Dictionaries["example"]
	if !ok {
		t.Errorf("Dictionary example failed to be loaded")
	}
	if len(dict.Definitions) != 82 {
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

package backend

import (
	"testing"
)

// Current test actually queries my anki database
func TestAnki(t *testing.T) {
	anki := NewAnkiInterface()
	note, err := anki.GetAnkiNote("将")
	if err != nil {
		t.Errorf(err.Error())
	}
	if note.Fields.Word != "将" || note.Fields.Sentence != "你将穿哪双鞋？" {
		t.Errorf("Failed to fetch correct card: %v", note)
	}
}

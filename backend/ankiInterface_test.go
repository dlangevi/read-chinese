package backend

import (
	"testing"
)

// Current test actually queries my anki database
func TestAnki(t *testing.T) {
	anki := NewAnkiInterface(testRuntime)
	err := anki.ImportAnkiReviewData()
	if err != nil {
		t.Error(err)
	}
	t.Error("ya")
	// note, err := anki.GetAnkiNote("将")
	//
	//	if err != nil {
	//		t.Errorf(err.Error())
	//	}
	//
	//	if note.Fields.Word != "将" || note.Fields.Sentence != "你将穿哪双鞋？" {
	//		t.Errorf("Failed to fetch correct card: %v", note)
	//	}
	//
	// apinote, resterr := anki.anki.Cards.Get("Hanzi:将")
	//
	//	if resterr != nil {
	//		t.Errorf("%v: %v", resterr.Error, resterr.Message)
	//	} else {
	//
	//		thenote := (*apinote)[0]
	//		if thenote.Interval < 100 {
	//			t.Errorf("Wrong interval, %v", thenote.Interval)
	//		}
	//	}
}

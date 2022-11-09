package backend

import (
	"os"
	"path"
	"testing"
)

func TestSegment(t *testing.T) {

	tmpMetaData := path.Join(os.TempDir(), "metadatabonk.json")
	LoadMetadata(tmpMetaData)
	defer os.Remove(tmpMetaData)
	d := NewDictionaries()
	d.AddDictionary("test", "testdata/testdata.json", "english")
	s, err := NewSegmentation(d)
	if err != nil {
		t.Errorf("Failed to load jieba %v", err)
	}
	sentence := "你好,我完全善良你。"
	cut := s.SegmentSentence(sentence)
	if len(cut) != 7 {
		t.Errorf("Basic test for now %v", cut)
	}

	sentence = "完全善良"
	cut = s.SegmentSentence(sentence)
	if len(cut) != 2 {
		t.Errorf("Basic test for now %v", cut)
	}
}

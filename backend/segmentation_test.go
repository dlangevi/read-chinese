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
	s, err := NewSegmentation(d)
	if err != nil {
		t.Errorf("Failed to load jieba %v", err)
	}
	sentence := "你好,我是david！我很高兴“认识”你。"
	cut := s.SegmentSentence(sentence)
	if len(cut) != 21 {
		t.Errorf("Basic test for now %v", cut)
	}

}

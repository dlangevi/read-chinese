package segmentation

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	x := m.Run()

	os.Exit(x)
}

func TestSegment(t *testing.T) {
	err := loadJieba()
	if err != nil {
		t.Errorf("Failed to load jieba %v", err)
	}
	sentence := "你好,我是david！我很高兴“认识”你。"
	cut := segmentSentence(sentence)
	if len(cut) != 18 {
		t.Errorf("Basic test for now %v", cut)
	}

}

package backend

import (
	"testing"
)

func TestPinyin(t *testing.T) {
	for _, example := range []struct {
		test   string
		expect string
	}{
		{"dong1", "dōng"},
		{"bao4 dao4", "bàodào"},
		{"xiang4 zheng1", "xiàngzhēng"},
		{"ban4 zhi2 min2 di4", "bànzhímíndì"},
		{"xiao4 lu:4", "xiàolǜ"},
	} {
		replace := ToStandardPinyin(example.test)
		if replace != example.expect {
			t.Error("Failed example", example, "got", replace)
		}
	}
}

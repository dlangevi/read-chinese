package backend

import (
	// "log"
	// "os"
	// "path"
	"strings"
	"testing"
)

func TestTextToSpeech(t *testing.T) {
	//  tts := NewTextToSpeach()
	// tmpMetaData := path.Join(os.TempDir(), "metadatabonk.json")
	//  LoadMetadata(tmpMetaData)
	//
	//  userSettings.AzureApiKey = "secret"
	//
	//  speechData, err := tts.Synthesize("我已经大四了，在就快毕业的时候我发现大学我什么也没办到")
	//  if err != nil {
	//    t.Errorf("Failed search %v", err)
	//  }
	//  if len(speechData) < 1000 {
	//    t.Errorf("Wrong search %v", speechData)
	//  }

	// tmpMetaData := path.Join(os.TempDir(), "metadatabonk.json")
	// userSettings, _ := LoadMetadata(tmpMetaData, testRuntime)
	// userSettings.AnkiConfig.AzureApiKey = ""
	// tts := NewTextToSpeach(userSettings)
	// voices, _ := tts.GetVoices()
	//  for _, voice := range voices {
	//    log.Println(
	//      voice.DisplayName,
	//      voice.Locale,
	//      voice.VoiceType,
	//      len(voice.StyleList),
	//      len(voice.RolePlayList))
	//
	//  }
	//  t.Errorf("Voices %v", len(voices))
}

func stripWhitespace(str string) string {
	return strings.Join(strings.Fields(str), "")
}

func compareRequests(t *testing.T, outputStr string, testStr string) {
	if stripWhitespace(outputStr) != stripWhitespace(testStr) {
		t.Errorf("Failed Format %s", outputStr)
	}
}

func TestSsmlGeneration(t *testing.T) {
	voice := Voice{
		Locale:        "",
		Voice:         "zh-CN-YunxiNeural",
		SpeakingStyle: "",
		RolePlay:      "",
		Speed:         0,
		Pitch:         0,
	}

	compareRequests(t, createSSMLRequest(voice, "你好"), `<speak
    xmlns="http://www.w3.org/2001/10/synthesis"
    xmlns:mstts="http://www.w3.org/2001/mstts"
    xmlns:emo="http://www.w3.org/2009/10/emotionml"
    version="1.0"
    xml:lang="en-US">
    <voice name="zh-CN-YunxiNeural">
      <prosody rate="0%" pitch="0%">
        你好
      </prosody>
    </voice>
  </speak>`)

	voice.SpeakingStyle = "Angry"
	compareRequests(t, createSSMLRequest(voice, "你好"), `<speak
    xmlns="http://www.w3.org/2001/10/synthesis"
    xmlns:mstts="http://www.w3.org/2001/mstts"
    xmlns:emo="http://www.w3.org/2009/10/emotionml"
    version="1.0"
    xml:lang="en-US">
    <voice name="zh-CN-YunxiNeural">
      <mstts:express-as style="Angry">
        <prosody rate="0%" pitch="0%">
          你好
        </prosody>
      </mstts:express-as>
    </voice>
  </speak>`)

	voice.RolePlay = "Boy"
	compareRequests(t, createSSMLRequest(voice, "你好"), `<speak
    xmlns="http://www.w3.org/2001/10/synthesis"
    xmlns:mstts="http://www.w3.org/2001/mstts"
    xmlns:emo="http://www.w3.org/2009/10/emotionml"
    version="1.0"
    xml:lang="en-US">
    <voice name="zh-CN-YunxiNeural">
      <mstts:express-as style="Angry" role="Boy">
        <prosody rate="0%" pitch="0%">
          你好
        </prosody>
      </mstts:express-as>
    </voice>
  </speak>`)

	voice.SpeakingStyle = ""
	compareRequests(t, createSSMLRequest(voice, "你好"), `<speak
    xmlns="http://www.w3.org/2001/10/synthesis"
    xmlns:mstts="http://www.w3.org/2001/mstts"
    xmlns:emo="http://www.w3.org/2009/10/emotionml"
    version="1.0"
    xml:lang="en-US">
    <voice name="zh-CN-YunxiNeural">
      <mstts:express-as role="Boy">
        <prosody rate="0%" pitch="0%">
          你好
        </prosody>
      </mstts:express-as>
    </voice>
  </speak>`)

	voice.Speed = 47
	voice.Pitch = -32
	compareRequests(t, createSSMLRequest(voice, "你好"), `<speak
    xmlns="http://www.w3.org/2001/10/synthesis"
    xmlns:mstts="http://www.w3.org/2001/mstts"
    xmlns:emo="http://www.w3.org/2009/10/emotionml"
    version="1.0"
    xml:lang="en-US">
    <voice name="zh-CN-YunxiNeural">
      <mstts:express-as role="Boy">
        <prosody rate="47%" pitch="-32%">
          你好
        </prosody>
      </mstts:express-as>
    </voice>
  </speak>`)

}

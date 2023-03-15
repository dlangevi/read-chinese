package backend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	REGION   = "eastus"
	ENDPOINT = "https://eastus.tts.speech.microsoft.com/"
	TTS      = "cognitiveservices/v1"
	VOICES   = "cognitiveservices/voices/list"
)

type TTSVoice struct {
	Name                string            `json:"Name"`
	DisplayName         string            `json:"DisplayName"`
	LocalName           string            `json:"LocalName"`
	ShortName           string            `json:"ShortName"`
	Gender              string            `json:"Gender"`
	Locale              string            `json:"Locale"`
	LocaleName          string            `json:"LocaleName"`
	StyleList           []string          `json:"StyleList,omitempty"`
	SecondaryLocaleList []string          `json:"SecondaryLocaleList,omitempty"`
	SampleRateHertz     string            `json:"SampleRateHertz"`
	VoiceType           string            `json:"VoiceType"`
	Status              string            `json:"Status"`
	ExtendedPropertyMap map[string]string `json:"ExtendedPropertyMap,omitempty"`
	RolePlayList        []string          `json:"RolePlayList,omitempty"`
	WordsPerMinute      string            `json:"WordsPerMinute"`
}

type TextToSpeech struct {
	httpClient   *resty.Client
	currentVoice int
	voices       []string
	userSettings *UserSettings
}

type Voice struct {
	Locale        string `json:"Locale"`
	Voice         string `json:"Voice"`
	SpeakingStyle string `json:"SpeakingStyle"`
	RolePlay      string `json:"RolePlay"`
	// -100 to 200
	Speed int `json:"Speed"`
	// -50 to 50
	Pitch int `json:"Pitch"`
}

func NewTextToSpeach(
	userSettings *UserSettings,
) *TextToSpeech {
	client := &TextToSpeech{
		httpClient:   resty.New(),
		currentVoice: 0,
		voices: []string{
			"zh-CN-YunxiNeural",
			"zh-CN-XiaochenNeural",
			"zh-CN-XiaoshuangNeural", // child
		},
		userSettings: userSettings,
	}
	client.httpClient.SetBaseURL(ENDPOINT)
	return client
}
func (tts *TextToSpeech) ExportVoice() TTSVoice {
	return TTSVoice{}
}

func (tts *TextToSpeech) GetVoices() ([]TTSVoice, error) {
	key := tts.userSettings.AzureConfig.AzureApiKey
	rsp, err := tts.httpClient.NewRequest().
		SetHeader("Content-Type", "application/ssml+xml").
		SetHeader("Ocp-Apim-Subscription-Key", key).
		Get(VOICES)
	if err != nil {
		return nil, err
	}

	voices := []TTSVoice{}
	err = json.Unmarshal(rsp.Body(), &voices)
	filteredVoices := []TTSVoice{}
	for _, voice := range voices {
		if strings.Contains(voice.Locale, "zh") {
			filteredVoices = append(filteredVoices, voice)
		}
	}
	return filteredVoices, err
}

func createSSMLRequest(voice Voice, text string) string {
	// Base of the request
	requestText := fmt.Sprintf(`<prosody rate="%d%%" pitch="%d%%">
      %v
    </prosody>`, voice.Speed, voice.Pitch, text)

	if voice.SpeakingStyle != "" && voice.RolePlay != "" {
		requestText = fmt.Sprintf(`<mstts:express-as style="%s" role="%s">
        %s
      </mstts:express-as>`,
			voice.SpeakingStyle, voice.RolePlay, requestText)
	} else if voice.SpeakingStyle != "" {
		requestText = fmt.Sprintf(`
      <mstts:express-as style="%s">
        %s
      </mstts:express-as>`, voice.SpeakingStyle, requestText)

	} else if voice.RolePlay != "" {
		requestText = fmt.Sprintf(`
      <mstts:express-as role="%s">
        %s
      </mstts:express-as>`, voice.RolePlay, requestText)
	}

	requestText = fmt.Sprintf(`<speak 
  xmlns="http://www.w3.org/2001/10/synthesis" 
  xmlns:mstts="http://www.w3.org/2001/mstts" 
  xmlns:emo="http://www.w3.org/2009/10/emotionml" 
  version="1.0" 
  xml:lang="en-US">
  <voice name="%s">
    %s
  </voice>
</speak>`, voice.Voice, requestText)

	return requestText
}

// We will use the rest api because I do not want to setup all the cgo sdk
func (tts *TextToSpeech) Synthesize(text string) (string, error) {
	voiceList := tts.userSettings.AzureConfig.VoiceList
	voice := voiceList[tts.currentVoice%len(voiceList)]
	tts.currentVoice += 1
	requestText := createSSMLRequest(voice, text)
	key := tts.userSettings.AzureConfig.AzureApiKey
	rsp, err := tts.httpClient.NewRequest().
		SetHeader("Content-Type", "application/ssml+xml").
		SetHeader("Ocp-Apim-Subscription-Key", key).
		SetHeader("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm").
		SetBody(requestText).
		Post(TTS)

	// Again, need to check for error somehow here
	body := rsp.Body()
	enc := base64.StdEncoding.EncodeToString(body)
	return enc, err
}

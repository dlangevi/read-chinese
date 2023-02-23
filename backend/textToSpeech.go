package backend

import (
	"encoding/base64"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	REGION = "eastus"
	URI    = "https://eastus.tts.speech.microsoft.com/cognitiveservices/v1"
)

type TextToSpeech struct {
	httpClient   *resty.Client
	currentVoice int
	voices       []string
	userSettings *UserConfig
}

func NewTextToSpeach(
	userSettings *UserConfig,
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
	client.httpClient.SetBaseURL(URI)
	return client
}

// We will use the rest api because I do not want to setup all the cgo sdk
func (tts *TextToSpeech) Synthesize(text string) (string, error) {
	voice := tts.voices[tts.currentVoice]
	tts.currentVoice = (tts.currentVoice + 1) % 3
	requestText := fmt.Sprintf(`<speak version='1.0' xml:lang='zh-CN'>
    <voice xml:lang='zh-CN' name='%v'>
      %v
    </voice>
  </speak>`, voice, text)

	key := tts.userSettings.AnkiConfig.AzureApiKey
	rsp, err := tts.httpClient.NewRequest().
		SetHeader("Content-Type", "application/ssml+xml").
		SetHeader("Ocp-Apim-Subscription-Key", key).
		SetHeader("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm").
		SetBody(requestText).
		Post("")

	// Again, need to check for error somehow here
	body := rsp.Body()
	enc := base64.StdEncoding.EncodeToString(body)
	return enc, err
}

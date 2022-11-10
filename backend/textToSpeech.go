package backend

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	REGION = "eastus"
	URI    = "https://eastus.tts.speech.microsoft.com/cognitiveservices/v1"
)

// TODO rotate voices
// const myVoices = [
//   'zh-CN-YunxiNeural',
//   'zh-CN-XiaochenNeural',
//   'zh-CN-XiaoshuangNeural', // child
// ];
// const nextVoice = (function nextVoice() {
//   let next = 0;
//   return () => {
//     next = (next + 1) % 3;
//     return next;
//   };
// }());

// We will use the rest api because I do not want to setup all the cgo sdk
func Synthesize(text string) (string, error) {
	client := &http.Client{}

	voice := "zh-CN-YunxiNeural"
	requestText := fmt.Sprintf(`<speak version='1.0' xml:lang='zh-CN'>
    <voice xml:lang='zh-CN' name='%v'>
      %v
    </voice>
  </speak>`, voice, text)

	req, err := http.NewRequest("POST", URI, bytes.NewBufferString(requestText))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/ssml+xml")
	key := userSettings.AzureApiKey
	req.Header.Set("Ocp-Apim-Subscription-Key", key)
	req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(rsp.Body)
		enc := base64.StdEncoding.EncodeToString(body)
		return enc, err
	} else {
		return "", errors.New("cannot convert text to speech")
	}
}

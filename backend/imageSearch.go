package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const (
	IMAGEURI = "https://api.bing.microsoft.com/v7.0/images/search?"
)

type ImageInfo struct {
	Name         string
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type ImageResponse struct {
	Value []ImageInfo
}

type ImageClient struct {
}

func (i *ImageClient) SearchImages(query string) ([]ImageInfo, error) {

	client := &http.Client{}
	// TODO URI encode
	params := url.Values{}
	params.Add("q", query)
	params.Add("count", "5")
	params.Add("imageType", "Photo")
	params.Add("safeSearch", "Strict")

	req, err := http.NewRequest("GET", IMAGEURI+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	key := userSettings.AzureImageApiKey
	req.Header.Set("Ocp-Apim-Subscription-Key", key)
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusOK {
		dec := json.NewDecoder(rsp.Body)
		response := &ImageResponse{}
		err := dec.Decode(&response)
		if err != nil {
			return nil, err
		}
		// TODO base64 encode the body which can then be sent directly
		// to anki connect?
		return response.Value, nil
	} else {
		return nil, errors.New("do images no gogo")
	}

}

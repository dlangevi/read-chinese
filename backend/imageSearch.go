package backend

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	IMAGEURI = "https://api.bing.microsoft.com"
)

type ImageInfo struct {
	Name         string
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type ImageResponse struct {
	Value []ImageInfo
}

type ImageClient struct {
	httpClient   *resty.Client
	userSettings *UserSettings
}

func NewImageClient(
	userSettings *UserSettings,
) *ImageClient {
	client := &ImageClient{
		httpClient:   resty.New(),
		userSettings: userSettings,
	}
	return client
}

func (i *ImageClient) SearchImages(query string) ([]ImageInfo, error) {
	i.httpClient.SetBaseURL(IMAGEURI)
	// i.httpClient.SetDisableWarn(true)
	result := &ImageResponse{}
	rsp, err := i.httpClient.R().
		SetHeader("Ocp-Apim-Subscription-Key", i.userSettings.AzureImageApiKey).
		SetQueryParams(map[string]string{
			"q":          query,
			"count":      "6",
			"imageType":  "Photo",
			"safeSearch": "Strict",
		}).
		SetResult(result).
		Get("/v7.0/images/search")
	// TODO need some better way of error detection
	if len(result.Value) != 6 {
		err = errors.New(fmt.Sprintf("Something failed with imageSearch %v", rsp))
	}
	return result.Value, err

}

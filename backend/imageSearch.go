package backend

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	IMAGEURI = "https://api.bing.microsoft.com"
)

// The format Azure responds in
type AzureImageInfo struct {
	Name         string
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type ImageInfo struct {
	Name      string `json:"name,omitempty"`
	Url       string `json:"url,omitempty"`
	ImageData string `json:"imageData,omitempty"`
}

type ImageResponse struct {
	Value []AzureImageInfo
}

type ImageClient struct {
	httpClient   *resty.Client
	userSettings *UserConfig
}

func NewImageClient(
	userSettings *UserConfig,
) *ImageClient {
	client := &ImageClient{
		httpClient:   resty.New(),
		userSettings: userSettings,
	}
	return client
}

func (i *ImageClient) SearchImages(query string) ([]ImageInfo, error) {
	i.httpClient.SetBaseURL(IMAGEURI)
	result := &ImageResponse{}
	rsp, err := i.httpClient.R().
		SetHeader("Ocp-Apim-Subscription-Key", i.userSettings.AnkiConfig.AzureImageApiKey).
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

	imageInfo := []ImageInfo{}
	for _, image := range result.Value {
		imageInfo = append(imageInfo, ImageInfo{
			Name: image.Name,
			Url:  image.ThumbnailUrl,
		})
	}

	return imageInfo, err

}

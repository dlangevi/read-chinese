package backend

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	IMAGEURI = "https://api.bing.microsoft.com"
)

type ThumbnailInfo struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

// The format Azure responds in
type AzureImageInfo struct {
	Name          string
	ThumbnailUrl  string
	ThumbnailInfo ThumbnailInfo `json:"thumbnail"`
}

type ImageInfo struct {
	Name        string `json:"name,omitempty"`
	Url         string `json:"url,omitempty"`
	ImageData   string `json:"imageData,omitempty"`
	ImageWidth  int64  `json:"imageWidth,omitEmpty"`
	ImageHeight int64  `json:"imageHeight,omitEmpty"`
}

type ImageResponse struct {
	Value []AzureImageInfo
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
	result := &ImageResponse{}
	rsp, err := i.httpClient.R().
		SetHeader("Ocp-Apim-Subscription-Key",
			i.userSettings.AzureConfig.AzureImageApiKey).
		// TODO look into region filtering
		// SetHeader("Accept-Language", "zh").
		SetQueryParams(map[string]string{
			"q":          query,
			"count":      "30",
			"safeSearch": "Strict",
			"imageType":  "Photo",
			"setLang":    "zh-hans",
			"cc":         "zh-CN",
			// "mkt":        "zh-CN",

		}).
		SetResult(result).
		Get("/v7.0/images/search")
		// TODO need some better way of error detection
	if len(result.Value) != 30 {
		err = errors.New(fmt.Sprintf("Something failed with imageSearch %v", rsp))
	}

	imageInfo := []ImageInfo{}
	for _, image := range result.Value {
		imageInfo = append(imageInfo, ImageInfo{
			Name:        image.Name,
			Url:         image.ThumbnailUrl,
			ImageWidth:  image.ThumbnailInfo.Width,
			ImageHeight: image.ThumbnailInfo.Height,
		})
	}

	return imageInfo, err

}

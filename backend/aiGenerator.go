package backend

import (
// some chatgpt go client
// some imagegen ai go client
)

type (
	AiGenerator interface {
		GenerateImage(description string) ImageInfo
		GenerateSentences(word string) []string
	}

	aiGenerator struct {
		// The api clients will probably go here
		userSettings *UserSettings
	}
)

func NewAiGenerator(
	userSettings *UserSettings,
) *aiGenerator {
	// Here setup whatever credentials (api keys etc) are needed
	// to setup whatever api clients are needed here
	return &aiGenerator{
		userSettings: userSettings,
	}
}

func (ai *aiGenerator) GenerateImage(description string) ImageInfo {
	// Call some image generation api. If it returns a url to the image,
	// Put that in the 'Url' field of imageinfo. Otherwise we need to
	// convert the image to a base64 string

	return ImageInfo{
		ImageData: "base64 image data here",
	}

}

func (ai *aiGenerator) GenerateSentences(word string) []string {
	// Call some chatgpt api with a promt like
	// "Generate 5 example sentences for the word ${word}
	// for a chinese language learner"
	sentences := []string{}
	return sentences
}

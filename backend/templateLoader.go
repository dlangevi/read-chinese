package backend

import (
	"embed"
	"encoding/json"
)

//go:embed assets/anki
var templates embed.FS

type (
	// Can interpret basic anki card templates and create an html
	// document that can then be displayed in the frontend to preview
	// various card types
	//
	// For now just support loading the default Card from assets, but in the
	// future the data can all be fetched from the users Anki
	TemplateLoader interface {
		GetTemplate() (AnkiTemplate, error)
	}

	CardFieldsMap struct {
		Hanzi               string
		Pinyin              string
		ChineseDefinition   string
		EnglishDefinition   string
		HanziAudio          string
		ExampleSentence     string
		Images              string
		SentenceTranslation string
		SentenceAudio       string
		Notes               string
	}

	AnkiTemplate struct {
		Name   string
		Front  string
		Back   string
		Css    string
		Fields map[string]string
	}

	templateLoader struct {
	}
)

func NewTemplateLoader() *templateLoader {
	return &templateLoader{}
}

func (c *templateLoader) GetTemplate() (AnkiTemplate, error) {
	template := AnkiTemplate{}

	fileBytes, err := templates.ReadFile(
		"assets/anki/templates/default/CardFront.html")
	if err != nil {
		return template, err
	}
	template.Front = string(fileBytes)
	fileBytes, err = templates.ReadFile(
		"assets/anki/templates/default/CardBack.html")
	if err != nil {
		return template, err
	}
	template.Back = string(fileBytes)
	fileBytes, err = templates.ReadFile(
		"assets/anki/templates/default/CardStyle.css")
	if err != nil {
		return template, err
	}
	template.Css = string(fileBytes)
	fileBytes, err = templates.ReadFile(
		"assets/anki/templates/default/CardFields.json")
	if err != nil {
		return template, err
	}
	err = json.Unmarshal(fileBytes, &template.Fields)
	if err != nil {
		return template, err
	}

	return template, nil
}

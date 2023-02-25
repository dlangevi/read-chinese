package backend

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/atselvan/ankiconnect"
	restError "github.com/privatesquare/bkst-go-utils/utils/errors"
	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

type Fields struct {
	Word        string `json:"word"`
	Sentence    string `json:"sentence"`
	EnglishDefn string `json:"englishDefn"`
	ChineseDefn string `json:"chineseDefn"`
	Pinyin      string `json:"pinyin"`
	// TODO merge these in some nice way
	ImageUrls []string `json:"imageUrls"`
	Image64   []string `json:"image64"`
}

type (
	AnkiInterface interface {
		GetAnkiNoteSkeleton(word string) RawAnkiNote
		CreateAnkiNote(fields Fields, tags []string) error
		GetAnkiNote(word string) (RawAnkiNote, error)
		UpdateNoteFields(noteID int64, fields Fields) (string, error)
		ImportAnkiKeywords() error
		LoadProblemCards() ([]ProblemCard, error)
		HealthCheck() string
		ConfigurationCheck() (string, error)
		LoadTemplate() error
		LoadModels() ([]string, error)
		LoadModelFields(model string) ([]string, error)
	}
)

type ankiInterface struct {
	anki         *ankiconnect.Client
	textToSpeech *TextToSpeech
	userSettings *UserConfig
	known        *KnownWords
}

func NewAnkiInterface(userSettings *UserConfig, known *KnownWords) *ankiInterface {
	return &ankiInterface{
		anki:         ankiconnect.NewClient(),
		textToSpeech: NewTextToSpeach(userSettings),
		userSettings: userSettings,
		known:        known,
	}
}

type RawAnkiNote struct {
	NoteId int64  `json:"noteId"`
	Fields Fields `json:"fields"`
}

func (a *ankiInterface) GetAnkiNoteSkeleton(word string) RawAnkiNote {
	return RawAnkiNote{
		Fields: Fields{
			Word: word,
		},
	}
}

func (a *ankiInterface) HealthCheck() string {
	restErr := a.anki.Ping()
	if restErr != nil {
		return "Could not connect to Anki"
	}
	return ""
}

func (a *ankiInterface) ConfigurationCheck() (string, error) {
	// ActiveDeck needs to be a real deck
	decks, restErr := a.anki.Decks.GetAll()
	if restErr != nil {
		return "", toError(restErr)
	}
	if !slices.Contains(*decks, a.userSettings.AnkiConfig.ActiveDeck) {
		return "Active Deck does exist in Anki", nil
	}

	// ActiveModel needs to be a real model
	models, restErr := a.anki.Models.GetAll()
	if restErr != nil {
		return "", toError(restErr)
	}
	activeModel := a.userSettings.AnkiConfig.ActiveModel
	if !slices.Contains(*models, activeModel) {
		return "Chose Note type does exist in Anki", nil
	}

	// CardConfiguration needs to exists and have certian fields
	// and those fields all need to be real fields
	currentModelFields, restErr := a.anki.Models.GetFields(activeModel)
	if restErr != nil {
		return "", toError(restErr)
	}
	fieldsConfig := a.userSettings.GetMapping(activeModel)

	allEmpty := true
	for fieldName, value := range map[string]string{
		"FirstField":        fieldsConfig.FirstField,
		"Hanzi":             fieldsConfig.Hanzi,
		"ExampleSentence":   fieldsConfig.ExampleSentence,
		"EnglishDefinition": fieldsConfig.EnglishDefinition,
		"ChineseDefinition": fieldsConfig.ChineseDefinition,
		"Pinyin":            fieldsConfig.Pinyin,
		"HanziAudio":        fieldsConfig.HanziAudio,
		"SentenceAudio":     fieldsConfig.SentenceAudio,
		"Images":            fieldsConfig.Images,
		"Notes":             fieldsConfig.Notes,
	} {
		if value == "" {
			continue
		} else if !slices.Contains(*currentModelFields, value) {
			return fmt.Sprintf(
				"Configured field for %v : (%v) does not exist",
				fieldName, value), nil
		} else {
			allEmpty = false
		}
	}

	if allEmpty {
		return "No fields have been configured for current model", nil
	}

	return "", nil
}

func (a *ankiInterface) getConfiguredMapping() (FieldsMapping, error) {

	currentModel := a.userSettings.AnkiConfig.ActiveModel
	currentSettings := a.userSettings.GetMapping(currentModel)
	// Required fields are:
	// Hanzi: yes
	// ExampleSentence: yes
	// EnglishDefinition <- one of these
	// ChineseDefinition <- one of these
	// HanziAudio: if genhanziaudio is set
	// SentenceAudio: if genAudio is set
	// Images: if imageApi is set
	// Notes: nope
	// For now just make sure hanzi is set
	if currentSettings.Hanzi == "" {
		return currentSettings, errors.New("Hanzi is not mapped")
	}

	return currentSettings, nil

}

func (a *ankiInterface) CreateAnkiNote(fields Fields, tags []string) error {
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return err
	}

	restErr := a.anki.Ping()
	if restErr != nil {
		return toError(restErr)
	}

	audio := []ankiconnect.Audio{}
	pictures := []ankiconnect.Picture{}

	addAudio := func(field string, dest string) error {
		audio64, err := a.textToSpeech.Synthesize(field)
		if err != nil {
			return err
		}
		audio = append(audio, ankiconnect.Audio{
			Data: audio64,
			Filename: fmt.Sprintf("read-chinese-%v-%v.wav",
				dest,
				time.Now().UnixMilli()),
			Fields: []string{
				dest,
			},
		})
		return nil
	}
	if a.userSettings.AnkiConfig.GenerateTermAudio {

		err := addAudio(fields.Word, currentMapping.HanziAudio)
		if err != nil {
			return err
		}
	}
	// dont generate if sentence is empty
	if a.userSettings.AnkiConfig.GenerateSentenceAudio && len(fields.Sentence) > 0 {
		err := addAudio(fields.Sentence, currentMapping.SentenceAudio)
		if err != nil {
			return err
		}
	}
	for i, image := range fields.ImageUrls {
		milli := time.Now().UnixMilli()
		pictures = append(pictures, ankiconnect.Picture{
			URL: image,
			// TODO dont guess the encoding format
			Filename: fmt.Sprintf("read-chinese-image-%v-%v.jpg", milli, i),
			Fields: []string{
				currentMapping.Images,
			},
		})
	}
	noteFields := ankiconnect.Fields{}

	if currentMapping.Hanzi != "" {
		noteFields[currentMapping.Hanzi] = fields.Word
	}

	if currentMapping.ExampleSentence != "" {
		noteFields[currentMapping.ExampleSentence] = fields.Sentence
	}

	if currentMapping.EnglishDefinition != "" {
		noteFields[currentMapping.EnglishDefinition] = fields.EnglishDefn
	}

	if currentMapping.ChineseDefinition != "" {
		noteFields[currentMapping.ChineseDefinition] = fields.ChineseDefn
	}

	if currentMapping.Pinyin != "" {
		noteFields[currentMapping.Pinyin] = fields.Pinyin
	}

	firstField := currentMapping.FirstField
	// If your first field is something else then ???
	if firstField != currentMapping.Hanzi ||
		firstField != currentMapping.ExampleSentence ||
		firstField != currentMapping.EnglishDefinition ||
		firstField != currentMapping.ChineseDefinition {
		noteFields[currentMapping.FirstField] = fmt.Sprint(time.Now().UnixMilli())
	}

	note := ankiconnect.Note{
		DeckName:  a.userSettings.AnkiConfig.ActiveDeck,
		ModelName: a.userSettings.AnkiConfig.ActiveModel,
		Fields:    noteFields,
		Tags:      tags,
		Options: &ankiconnect.Options{
			AllowDuplicate: a.userSettings.AnkiConfig.AllowDuplicates,
		},
		Audio:   audio,
		Picture: pictures,
	}
	restErr = a.anki.Notes.Add(note)
	if restErr != nil {
		return toError(restErr)
	}

	a.known.AddWord(fields.Word, 0)

	return nil
}

func toError(restErr *restError.RestErr) error {
	return errors.New(fmt.Sprintf("%v, %v", restErr.Error, restErr.Message))
}

func (a *ankiInterface) GetAnkiNote(word string) (RawAnkiNote, error) {
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return RawAnkiNote{}, err
	}

	notes, restErr := a.anki.Notes.Get(fmt.Sprintf("%v:%v", currentMapping.Hanzi, word))
	if restErr != nil {
		return RawAnkiNote{}, toError(restErr)
	}
	if len(*notes) == 0 {
		return RawAnkiNote{}, errors.New("No note exists")
	}
	if len(*notes) > 1 {
		return RawAnkiNote{}, errors.New("Duplicate notes exists")
	}

	note := (*notes)[0]
	extract := func(field string) string {
		value, _ := note.Fields[field]
		return value.Value
	}

	images := []string{}
	// Images will (hopefully) have html <img src=\"filename.jpg\">
	imagesString := extract(currentMapping.Images)
	if imagesString != "" {
		log.Print(imagesString)
		doc, err := html.Parse(strings.NewReader(imagesString))
		if err != nil {
			return RawAnkiNote{}, err
		}
		imageNames := []string{}
		//
		var f func(*html.Node)
		// This is definitly a little overkill but hey
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "img" {
				for _, a := range n.Attr {
					if a.Key == "src" {
						imageNames = append(imageNames, a.Val)
						break
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		for _, name := range imageNames {
			image, err := a.anki.Media.RetrieveMediaFile(name)
			if err != nil {
				return RawAnkiNote{}, toError(err)
			}
			images = append(images, *image)
		}

	}

	rawNote := RawAnkiNote{
		NoteId: note.NoteId,
		Fields: Fields{
			Word:        extract(currentMapping.Hanzi),
			Sentence:    extract(currentMapping.ExampleSentence),
			EnglishDefn: extract(currentMapping.EnglishDefinition),
			ChineseDefn: extract(currentMapping.ChineseDefinition),
			Pinyin:      extract(currentMapping.Pinyin),
			// ImageUrls: [],
			Image64: images,
		},
	}
	return rawNote, nil
}

func (a *ankiInterface) UpdateNoteFields(noteID int64, fields Fields) (string, error) {
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return "", err
	}
	ankiFields := ankiconnect.Fields{}
	if fields.Word != "" {
		ankiFields[currentMapping.Hanzi] = fields.Word
	}
	if fields.Sentence != "" {
		ankiFields[currentMapping.ExampleSentence] = fields.Sentence
	}
	if fields.EnglishDefn != "" {
		ankiFields[currentMapping.EnglishDefinition] = fields.EnglishDefn
	}
	if fields.ChineseDefn != "" {
		ankiFields[currentMapping.ChineseDefinition] = fields.ChineseDefn
	}
	if fields.Pinyin != "" {
		ankiFields[currentMapping.Pinyin] = fields.Pinyin
	}
	// TODO include new audio and images
	update := ankiconnect.UpdateNote{
		Id:     noteID,
		Fields: ankiFields,
	}
	restErr := a.anki.Notes.Update(update)
	if restErr != nil {
		return "", toError(restErr)
	}

	// TODO clear flag if needed

	return "success", nil
}

func (a *ankiInterface) ImportAnkiKeywords() error {
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return err
	}
	cards, restErr := a.anki.Cards.Get(
		fmt.Sprintf(`"deck:%v"`, a.userSettings.AnkiConfig.ActiveDeck))
	if restErr != nil {
		return toError(restErr)
	}

	words := []WordEntry{}
	for _, card := range *cards {
		word, _ := card.Fields[currentMapping.Hanzi]
		words = append(words, WordEntry{
			Word:     word.Value,
			Interval: card.Interval,
		})
	}
	return a.known.AddWords(words)
}

type FlaggedCard struct {
	Word     string `json:"word"`
	Sentence string `json:"sentence"`
}

// Worry about using an enum or something later. I know this will work
// into Wails
type Problems struct {
	Flagged              bool `json:"Flagged"`
	MissingImage         bool `json:"MissingImage"`
	MissingSentence      bool `json:"MissingSentence"`
	MissingSentenceAudio bool `json:"MissingSentenceAudio"`
	MissingWordAudio     bool `json:"MissingWordAudio"`
	MissingPinyin        bool `json:"MissingPinyin"`
}

type ProblemCard struct {
	Word     string   `json:"Word"`
	Problems Problems `json:"Problems"`
	Notes    string   `json:"Notes"`
}

func (a *ankiInterface) LoadProblemCards() ([]ProblemCard, error) {
	problemCardsMap := map[int64]ProblemCard{}
	problemCards := []ProblemCard{}
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return problemCards, err
	}

	// For each query, get the noteIds. Add them all to the map
	// with problems

	type ProblemCase struct {
		Query  string
		Setter func(*Problems)
	}

	currentDeck := a.userSettings.AnkiConfig.ActiveDeck

	checks := []ProblemCase{
		{ // Cards Flagged by the user
			Query:  fmt.Sprintf("deck:%v -flag:0", currentDeck),
			Setter: func(p *Problems) { p.Flagged = true },
		},
		{ // Missing Example Sentence
			Query:  fmt.Sprintf("deck:%v %v:", currentDeck, currentMapping.ExampleSentence),
			Setter: func(p *Problems) { p.MissingSentence = true },
		},
		{ // Has Sentence, but missing Sentence Audio
			Query: fmt.Sprintf("deck:%v -%v: %v:", currentDeck,
				currentMapping.ExampleSentence, currentMapping.SentenceAudio),
			Setter: func(p *Problems) { p.MissingSentenceAudio = true },
		},
		{ // Missing Image
			Query:  fmt.Sprintf("deck:%v %v:", currentDeck, currentMapping.Images),
			Setter: func(p *Problems) { p.MissingImage = true },
		},
		{ // Missing HanziAudio
			Query:  fmt.Sprintf("deck:%v %v:", currentDeck, currentMapping.HanziAudio),
			Setter: func(p *Problems) { p.MissingWordAudio = true },
		},
		{ // Missing Pinyin TODO: check for ugly pinyin? eg: ni3hao3
			Query:  fmt.Sprintf("deck:%v %v:", currentDeck, currentMapping.Pinyin),
			Setter: func(p *Problems) { p.MissingPinyin = true },
		},
	}

	for _, check := range checks {
		// Todo switch out Get for Search + Lookup if we want to speed it up
		ids, restErr := a.anki.Notes.Get(check.Query)
		if restErr != nil {
			return nil, toError(restErr)
		}

		for _, id := range *ids {
			word, ok := id.Fields[currentMapping.Hanzi]
			if !ok {
				return nil, errors.New("Hanzi not found")
			}
			problemCard, exists := problemCardsMap[id.NoteId]
			if !exists {
				problemCard = ProblemCard{
					Word:     word.Value,
					Problems: Problems{},
				}

				notes, ok := id.Fields[currentMapping.Notes]
				if ok {
					problemCard.Notes = notes.Value
				}
			}
			check.Setter(&problemCard.Problems)
			problemCardsMap[id.NoteId] = problemCard
		}

	}

	// Big lookup on all selected Note ids to map to the word in the field
	// For now just do

	// Finally map to array since wails cant do the needed ts stuff
	for _, problemCard := range problemCardsMap {
		problemCards = append(problemCards, problemCard)
	}

	return problemCards, nil
}

func (a *ankiInterface) LoadTemplate() error {
	template, err := NewTemplateLoader().GetTemplate()
	if err != nil {
		return err
	}

	model := ankiconnect.Model{
		ModelName:     template.Name,
		InOrderFields: template.Fields,
		Css:           template.Css,
		IsCloze:       false,
		CardTemplates: []ankiconnect.CardTemplate{
			{
				Name:  template.Name,
				Front: template.Front,
				Back:  template.Back,
			},
		},
	}
	restErr := a.anki.Models.Create(model)
	if restErr != nil {
		return toError(restErr)
	}
	return nil
}

func (a *ankiInterface) LoadModels() ([]string, error) {
	models, restErr := a.anki.Models.GetAll()
	if restErr != nil {
		return nil, toError(restErr)
	}
	return *models, nil
}

func (a *ankiInterface) LoadDecks() ([]string, error) {
	decks, restErr := a.anki.Decks.GetAll()
	if restErr != nil {
		return nil, toError(restErr)
	}
	return *decks, nil
}

func (a *ankiInterface) LoadModelFields(model string) ([]string, error) {
	fields, restErr := a.anki.Models.GetFields(model)
	if restErr != nil {
		return nil, toError(restErr)
	}
	return *fields, nil

}

// TODO in client
// setSpecificValueOfCard
// and use that to unflag stuff when updated

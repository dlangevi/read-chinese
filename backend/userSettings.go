package backend

import (
	"errors"
	"log"
	"os"

	"encoding/json"
	"github.com/oleiade/reflections"
)

type Dict struct {
	Path     string
	Language string
}

type FieldsMapping struct {
	Hanzi             string `json:"hanzi"`
	ExampleSentence   string `json:"exampleSentence"`
	EnglishDefinition string `json:"englishDefinition"`
	ChineseDefinition string `json:"chineseDefinition"`
	Pinyin            string `json:"pinyin"`
	HanziAudio        string `json:"hanziAudio"`
	SentenceAudio     string `json:"sentenceAudio"`
	Images            string `json:"images"`
	Notes             string `json:"notes"`
}

type UserSettings struct {
	// Meta fields
	path string
	Ran  int

	// Card Creation Settings
	// This is what the user interacts with
	// in the front end
	AutoAdvanceSentence bool
	PopulateEnglish     bool
	PopulateChinese     bool
	AutoAdvanceEnglish  bool
	AutoAdvanceImage    bool
	AutoAdvanceCard     bool

	// Anki / Card generation
	// This controls how the front end gets
	// translated into Anki calls
	ActiveDeck            string
	ActiveModel           string
	ModelMappings         map[string]FieldsMapping
	AddProgramTag         bool
	AddBookTag            bool
	AllowDuplicates       bool
	GenerateTermAudio     bool
	GenerateSentenceAudio bool
	AzureApiKey           string
	AzureImageApiKey      string

	// Dictionaries
	Dicts           map[string]Dict
	PrimaryDict     string
	ShowDefinitions bool
	EnableChinese   bool

	// Sentence Generation
	IdealSentenceLength int
	KnownInterval       int

	// Book Library
	OnlyFavorites bool
	HideRead      bool

	// Card Management
	ProblemFlagged              bool
	ProblemMissingImage         bool
	ProblemMissingSentence      bool
	ProblemMissingSentenceAudio bool
	ProblemMissingWordAudio     bool
	ProblemMissingPinyin        bool
}

func defaultSettings(path string) *UserSettings {
	return &UserSettings{
		// Meta fields
		path: path,
		Ran:  0,

		// Card Creation
		AutoAdvanceSentence: true,
		PopulateEnglish:     false,
		PopulateChinese:     false,
		AutoAdvanceEnglish:  false,
		AutoAdvanceImage:    false,
		AutoAdvanceCard:     true,

		// Anki / Card generation
		ActiveDeck:  "Reading",
		ActiveModel: "Reading Card",
		ModelMappings: map[string]FieldsMapping{
			"Reading Card": {
				Hanzi:             "Hanzi",
				ExampleSentence:   "ExampleSentence",
				EnglishDefinition: "EnglishDefinition",
				ChineseDefinition: "ChineseDefinition",
				Pinyin:            "Pinyin",
				HanziAudio:        "HanziAudio",
				SentenceAudio:     "SentenceAudio",
				Images:            "Images",
				Notes:             "Notes",
			},
		},
		AddProgramTag:   true,
		AddBookTag:      true,
		AllowDuplicates: true,

		GenerateTermAudio:     false,
		GenerateSentenceAudio: false,
		AzureApiKey:           "",
		AzureImageApiKey:      "",

		// Dictionaries
		Dicts:           map[string]Dict{},
		PrimaryDict:     "",
		ShowDefinitions: true,
		EnableChinese:   true,

		// Sentence Generation
		KnownInterval:       10,
		IdealSentenceLength: 20,

		// Book Library
		OnlyFavorites: false,
		HideRead:      false,

		// Card Management
		ProblemFlagged:              true,
		ProblemMissingImage:         true,
		ProblemMissingSentence:      true,
		ProblemMissingSentenceAudio: true,
		ProblemMissingWordAudio:     true,
		ProblemMissingPinyin:        true,
	}
}

func LoadMetadata(path string) (*UserSettings, error) {
	userSettings := defaultSettings(path)

	if _, err := os.Stat(path); err == nil {
		// metadata already exists, read from it
		b, err := os.ReadFile(path)
		if err = json.Unmarshal(b, userSettings); err != nil {
			return nil, err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// metadata does *not* exist, write the default settings
		userSettings.saveMetadata()
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.Fatal(err)
	}
	return userSettings, nil
}

func (m *UserSettings) saveMetadata() error {
	path := m.path
	str, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, str, 0666); err != nil {
		return err
	}
	return nil
}

func (m *UserSettings) UpdateTimesRan() {
	m.Ran += 1
	m.saveMetadata()
}

func (m *UserSettings) GetTimesRan() int {
	return m.Ran
}

func (m *UserSettings) GetUserSetting(key string) string {
	return getUserSetting[string](m, key)
}

func (m *UserSettings) GetUserSettingBool(key string) bool {
	return getUserSetting[bool](m, key)
}

func (m *UserSettings) GetUserSettingInt(key string) int {
	return getUserSetting[int](m, key)
}

func getUserSetting[T int | string | bool](m *UserSettings, key string) T {
	val, err := reflections.GetField(m, key)
	if err != nil {
		log.Fatal(err)
	}
	return val.(T)
}

func (m *UserSettings) SetUserSetting(key string, val string) {
	setUserSetting(m, key, val)
	m.saveMetadata()
}

func (m *UserSettings) SetUserSettingBool(key string, val bool) {
	setUserSetting(m, key, val)
	m.saveMetadata()
}

func (m *UserSettings) SetUserSettingInt(key string, val int) {
	setUserSetting(m, key, val)
	m.saveMetadata()
}

func setUserSetting[T int | string | bool](m *UserSettings, key string, val T) {
	err := reflections.SetField(m, key, val)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *UserSettings) SaveDict(name string, dictPath string, language string) {
	m.Dicts[name] = Dict{
		Path:     dictPath,
		Language: language,
	}
	m.saveMetadata()
}

func (m *UserSettings) DeleteDict(name string) {
	delete(m.Dicts, name)
	m.saveMetadata()
}

func (m *UserSettings) SetPrimaryDict(dictName string) {
	// TODO Make sure its a real dict
	m.PrimaryDict = dictName
	m.saveMetadata()
}

func (m *UserSettings) ExportMapping() FieldsMapping {
	return FieldsMapping{}

}

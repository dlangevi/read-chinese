package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
)

type (
	UserConfig struct {
		// Meta fields
		Meta MetaSettings `json:"meta"`

		// Card Creation Settings
		CardCreationConfig CardCreationConfig `json:"CardCreation"`

		// Anki / Card generation
		AnkiConfig AnkiConfig `json:"AnkiConfig"`

		// Dictionaries
		DictionariesConfig DictionaryConfig `json:"Dictionaries"`

		// Sentence Generation
		SentenceGenerationConfig SentenceGenerationConfig `json:"SentenceGeneration"`

		// Book Library
		LibraryConfig LibraryConfig `json:"BookLibrary"`

		// Card Management
		CardManagementConfig CardManagementConfig `json:"CardManagement"`
	}

	MetaSettings struct {
		path string
		Ran  int `json:"Ran"`
	}

	CardCreationConfig struct {
		AutoAdvanceSentence bool `json:"AutoAdvanceSentence"`
		PopulateEnglish     bool `json:"PopulateEnglish"`
		PopulateChinese     bool `json:"PopulateChinese"`
		AutoAdvanceEnglish  bool `json:"AutoAdvanceEnglish"`
		AutoAdvanceImage    bool `json:"AutoAdvanceImage"`
		AutoAdvanceCard     bool `json:"AutoAdvanceCard"`
	}

	FieldsMapping struct {
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
	AnkiConfig struct {
		ActiveDeck            string                   `json:"ActiveDeck"`
		ActiveModel           string                   `json:"ActiveModel"`
		ModelMappings         map[string]FieldsMapping `json:"ModelMappings"`
		AddProgramTag         bool                     `json:"AddProgramTag"`
		AddBookTag            bool                     `json:"AddBookTag"`
		AllowDuplicates       bool                     `json:"AllowDuplicates"`
		GenerateTermAudio     bool                     `json:"GenerateTermAudio"`
		GenerateSentenceAudio bool                     `json:"GenerateSentenceAudio"`
		AzureApiKey           string                   `json:"AzureApiKey"`
		AzureImageApiKey      string                   `json:"AzureImageApiKey"`
	}

	Dict struct {
		Path     string `json:"Path"`
		Language string `json:"Language"`
	}

	DictionaryConfig struct {
		Dicts           map[string]Dict `json:"Dicts"`
		PrimaryDict     string          `json:"PrimaryDict"`
		ShowDefinitions bool            `json:"ShowDefinitions"`
		EnableChinese   bool            `json:"EnableChinese"`
	}

	SentenceGenerationConfig struct {
		IdealSentenceLength int `json:"IdealSentenceLength"`
		KnownInterval       int `json:"KnownInterval"`
	}

	LibraryConfig struct {
		OnlyFavorites bool `json:"OnlyFavorites"`
		HideRead      bool `json:"HideRead"`
	}

	CardManagementConfig struct {
		ProblemFlagged              bool `json:"ProblemFlagged"`
		ProblemMissingImage         bool `json:"ProblemMissingImage"`
		ProblemMissingSentence      bool `json:"ProblemMissingSentence"`
		ProblemMissingSentenceAudio bool `json:"ProblemMissingSentenceAudio"`
		ProblemMissingWordAudio     bool `json:"ProblemMissingWordAudio"`
		ProblemMissingPinyin        bool `json:"ProblemMissingPinyin"`
	}
)

func defaultConfig(path string) *UserConfig {
	return &UserConfig{
		Meta: MetaSettings{
			path: path,
			Ran:  0,
		},

		CardCreationConfig: CardCreationConfig{
			AutoAdvanceSentence: true,
			PopulateEnglish:     false,
			PopulateChinese:     false,
			AutoAdvanceEnglish:  false,
			AutoAdvanceImage:    false,
			AutoAdvanceCard:     true,
		},

		AnkiConfig: AnkiConfig{
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
			AddProgramTag:         true,
			AddBookTag:            true,
			AllowDuplicates:       true,
			GenerateTermAudio:     false,
			GenerateSentenceAudio: false,
			AzureApiKey:           "",
			AzureImageApiKey:      "",
		},

		DictionariesConfig: DictionaryConfig{
			Dicts:           map[string]Dict{},
			PrimaryDict:     "",
			ShowDefinitions: true,
			EnableChinese:   true,
		},

		SentenceGenerationConfig: SentenceGenerationConfig{
			KnownInterval:       10,
			IdealSentenceLength: 20,
		},

		LibraryConfig: LibraryConfig{
			OnlyFavorites: false,
			HideRead:      false,
		},

		CardManagementConfig: CardManagementConfig{
			ProblemFlagged:              true,
			ProblemMissingImage:         true,
			ProblemMissingSentence:      true,
			ProblemMissingSentenceAudio: true,
			ProblemMissingWordAudio:     true,
			ProblemMissingPinyin:        true,
		},
	}
}

func getValue(settings *UserConfig, field string) interface{} {
	value := reflect.ValueOf(settings).Elem()

	for i := 0; i < value.NumField(); i++ {
		subVal := value.Field(i)
		if subVal.Kind() != reflect.Struct {
			continue
		}

		subFieldVal := subVal.FieldByName(field)
		if subFieldVal.IsValid() {
			return subFieldVal.Interface()
		}
	}

	return nil
}

func setValue(settings *UserConfig, field string, newValue interface{}) error {
	value := reflect.ValueOf(settings).Elem()

	for i := 0; i < value.NumField(); i++ {
		subVal := value.Field(i)
		if subVal.Kind() != reflect.Struct {
			continue
		}

		subFieldVal := subVal.FieldByName(field)
		if subFieldVal.IsValid() {
			if !subFieldVal.CanSet() {
				return fmt.Errorf("cannot set field %s", field)
			}
			structFieldType := subFieldVal.Type()
			val := reflect.ValueOf(newValue)
			if val.Type().AssignableTo(structFieldType) {
				subFieldVal.Set(val)
				return nil
			} else {
				return fmt.Errorf("value type %s cannot be assigned to field %s", val.Type(), field)
			}
		}
	}

	return fmt.Errorf("field %s not found", field)
}

func LoadMetadata(path string) (*UserConfig, error) {
	userSettings := defaultConfig(path)

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

func (m *UserConfig) GetUserSettings() UserConfig {
	return *m
}

func (m *UserConfig) saveMetadata() error {
	path := m.Meta.path
	str, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, str, 0666); err != nil {
		return err
	}

	// runtime.EventsEmit(*m.ctx, "UpdatedConfig", m)
	return nil
}

func (m *UserConfig) UpdateTimesRan() {
	m.Meta.Ran += 1
	m.saveMetadata()
}

func (m *UserConfig) GetTimesRan() int {
	return m.Meta.Ran
}

func (m *UserConfig) SetUserSetting(key string, val string) {
	setUserSetting(m, key, val)
	m.saveMetadata()
}

func (m *UserConfig) SetUserSettingBool(key string, val bool) {
	setUserSetting(m, key, val)
	m.saveMetadata()
}

func (m *UserConfig) SetUserSettingInt(key string, val int) {
	setUserSetting(m, key, val)
	m.saveMetadata()
}

func setUserSetting[T int | string | bool](m *UserConfig, key string, val T) {
	err := setValue(m, key, val)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *UserConfig) SaveDict(name string, dictPath string, language string) {
	m.DictionariesConfig.Dicts[name] = Dict{
		Path:     dictPath,
		Language: language,
	}
	m.saveMetadata()
}

func (m *UserConfig) DeleteDict(name string) {
	delete(m.DictionariesConfig.Dicts, name)
	m.saveMetadata()
}

func (m *UserConfig) SetPrimaryDict(dictName string) {
	// TODO Make sure its a real dict
	m.DictionariesConfig.PrimaryDict = dictName
	m.saveMetadata()
}

func (m *UserConfig) GetMapping(modelName string) FieldsMapping {
	mapping, ok := m.AnkiConfig.ModelMappings[modelName]
	if !ok {
		return FieldsMapping{}
	}
	return mapping
}

func (m *UserConfig) SetMapping(modelName string, mapping FieldsMapping) error {
	m.AnkiConfig.ModelMappings[modelName] = mapping
	m.saveMetadata()
	return nil
}

func (m *UserConfig) DeleteMapping(modelName string) error {
	delete(m.AnkiConfig.ModelMappings, modelName)
	m.saveMetadata()
	return nil
}

func (m *UserConfig) ExportMapping() FieldsMapping {
	return FieldsMapping{}

}

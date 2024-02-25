package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	UserSettings struct {
		backend *Backend
		// Meta fields
		Meta MetaSettings `json:"meta"`

		// Card Creation Settings
		CardCreationConfig CardCreationConfig `json:"CardCreation"`

		// Anki / Card generation
		AnkiConfig AnkiConfig `json:"AnkiConfig"`

		// Azure Settings
		AzureConfig AzureConfig `json:"AzureConfig"`

		// Dictionaries
		DictionariesConfig DictionaryConfig `json:"Dictionaries"`

		WordLists WordListsConfig `json:"WordLists"`

		// Sentence Generation
		SentenceGenerationConfig SentenceGenerationConfig `json:"SentenceGeneration"`

		// Book Library
		LibraryConfig LibraryConfig `json:"BookLibrary"`

		// Card Management
	}

	MetaSettings struct {
		path               string
		EnableExperimental bool   `json:"EnableExperimental"`
		Ran                int    `json:"Ran"`
		Theme              string `json:"Theme"`
	}

	CardCreationConfig struct {
		AutoAdvanceSentence bool `json:"AutoAdvanceSentence"`
		PopulateEnglish     bool `json:"PopulateEnglish"`
		PopulateChinese     bool `json:"PopulateChinese"`
		AutoAdvanceEnglish  bool `json:"AutoAdvanceEnglish"`
		AutoAdvanceChinese  bool `json:"AutoAdvanceChinese"`
		AutoAdvanceImage    bool `json:"AutoAdvanceImage"`
		AutoAdvanceCard     bool `json:"AutoAdvanceCard"`
	}

	FieldsMapping struct {
		FirstField        string `json:"firstField"`
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
		ActiveDeck      string                   `json:"ActiveDeck"`
		ActiveModel     string                   `json:"ActiveModel"`
		ModelMappings   map[string]FieldsMapping `json:"ModelMappings"`
		AddProgramTag   bool                     `json:"AddProgramTag"`
		AddBookTag      bool                     `json:"AddBookTag"`
		AllowDuplicates bool                     `json:"AllowDuplicates"`
	}

	AzureConfig struct {
		GenerateTermAudio     bool    `json:"GenerateTermAudio"`
		GenerateSentenceAudio bool    `json:"GenerateSentenceAudio"`
		AzureApiKey           string  `json:"AzureApiKey"`
		AzureImageApiKey      string  `json:"AzureImageApiKey"`
		VoiceList             []Voice `json:"VoiceList"`
	}

	Dict struct {
		Path     string `json:"Path"`
		Language string `json:"Language"`
	}

	DictionaryConfig struct {
		Dicts           map[string]Dict `json:"Dicts"`
		PrimaryDict     string          `json:"PrimaryDict"`
		ShowDefinitions bool            `json:"ShowDefinitions"`
		ShowPinyin      bool            `json:"ShowPinyin"`
		EnableChinese   bool            `json:"EnableChinese"`
	}

	WordListsConfig struct {
		WordLists       map[string]string `json:"WordLists"`
		PrimaryWordList string            `json:"PrimaryWordList"`
	}

	SentenceGenerationConfig struct {
		IdealSentenceLength int `json:"IdealSentenceLength"`
		KnownInterval       int `json:"KnownInterval"`
	}

	LibraryConfig struct {
		OnlyFavorites bool `json:"OnlyFavorites"`
		HideRead      bool `json:"HideRead"`
		// If only go had proper enum
		DisplayTable bool `json:"DisplayTable"`
	}
)

func defaultConfig(path string) *UserSettings {
	return &UserSettings{
		Meta: MetaSettings{
			path:               path,
			EnableExperimental: false,
			Ran:                0,
			Theme:              "emerald",
		},

		CardCreationConfig: CardCreationConfig{
			AutoAdvanceSentence: false,
			PopulateEnglish:     false,
			PopulateChinese:     false,
			AutoAdvanceEnglish:  false,
			AutoAdvanceChinese:  false,
			AutoAdvanceImage:    false,
			AutoAdvanceCard:     false,
		},

		AnkiConfig: AnkiConfig{
			ActiveDeck:  "",
			ActiveModel: "",
			ModelMappings: map[string]FieldsMapping{
				"read-chinese-note": {
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
		},

		AzureConfig: AzureConfig{
			GenerateTermAudio:     false,
			GenerateSentenceAudio: false,
			AzureApiKey:           "",
			AzureImageApiKey:      "",
			VoiceList: []Voice{
				{
					Locale:   "zh-CN",
					Voice:    "zh-CN-YunxiNeural",
					RolePlay: "Narrator",
				},
				{
					Locale: "zh-CN",
					Voice:  "zh-CN-XiaochenNeural",
				},
				{
					Locale:        "zh-CN",
					Voice:         "zh-CN-YunxiNeural",
					SpeakingStyle: "narration-relaxed",
					RolePlay:      "Boy",
				},
				{
					Locale: "zh-CN",
					Voice:  "zh-CN-XiaoshuangNeural",
				},
			},
		},

		DictionariesConfig: DictionaryConfig{
			Dicts:           map[string]Dict{},
			PrimaryDict:     "",
			ShowDefinitions: true,
			ShowPinyin:      true,
			EnableChinese:   false,
		},

		WordLists: WordListsConfig{
			WordLists:       map[string]string{},
			PrimaryWordList: "",
		},

		SentenceGenerationConfig: SentenceGenerationConfig{
			KnownInterval:       10,
			IdealSentenceLength: 20,
		},

		LibraryConfig: LibraryConfig{
			OnlyFavorites: false,
			HideRead:      false,
			DisplayTable:  false,
		},
	}
}

func getValue(settings *UserSettings, field string) interface{} {
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

func setValue(settings *UserSettings, field string, newValue interface{}) error {
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

func LoadMetadata(path string, backend *Backend) (*UserSettings, error) {
	userSettings := defaultConfig(path)
	userSettings.backend = backend

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

func (m *UserSettings) GetUserSettings() UserSettings {
	return *m
}

func (m *UserSettings) saveMetadata() error {
	path := m.Meta.path
	str, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, str, 0666); err != nil {
		return err
	}

	if m.backend.ctx != nil {
		runtime.EventsEmit(m.backend.ctx, "UpdatedConfig", m)
	}
	return nil
}

func (m *UserSettings) UpdateTimesRan() {
	m.Meta.Ran += 1
	m.saveMetadata()
}

func (m *UserSettings) GetTimesRan() int {
	return m.Meta.Ran
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
	err := setValue(m, key, val)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *UserSettings) SaveDict(name string, dictPath string, language string) {
	m.DictionariesConfig.Dicts[name] = Dict{
		Path:     dictPath,
		Language: language,
	}
	m.saveMetadata()
}

func (m *UserSettings) DeleteDict(name string) {
	delete(m.DictionariesConfig.Dicts, name)
	m.saveMetadata()
}

func (m *UserSettings) SetPrimaryDict(dictName string) {
	// TODO Make sure its a real dict
	m.DictionariesConfig.PrimaryDict = dictName
	m.saveMetadata()
}

func (m *UserSettings) SaveList(name string, listPath string) {
	m.WordLists.WordLists[name] = listPath
	m.saveMetadata()
}

func (m *UserSettings) DeleteList(name string) {
	delete(m.WordLists.WordLists, name)
	m.saveMetadata()
}

func (m *UserSettings) SetPrimaryList(listName string) {
	// TODO Make sure its a real list
	m.WordLists.PrimaryWordList = listName
	m.saveMetadata()
}

func (m *UserSettings) GetMapping(modelName string) FieldsMapping {
	mapping, ok := m.AnkiConfig.ModelMappings[modelName]
	if !ok {
		return FieldsMapping{}
	}
	return mapping
}

func (m *UserSettings) SetMapping(modelName string, mapping FieldsMapping) error {
	m.AnkiConfig.ModelMappings[modelName] = mapping
	m.saveMetadata()
	return nil
}

func (m *UserSettings) DeleteMapping(modelName string) error {
	delete(m.AnkiConfig.ModelMappings, modelName)
	m.saveMetadata()
	return nil
}

func (m *UserSettings) AddVoice(voice Voice) error {
	m.AzureConfig.VoiceList = append(m.AzureConfig.VoiceList, voice)
	m.saveMetadata()
	return nil
}

func (m *UserSettings) RemoveVoice(voice Voice) error {
	voiceList := []Voice{}
	for _, currentVoice := range m.AzureConfig.VoiceList {
		if voice != currentVoice {
			voiceList = append(voiceList, currentVoice)
		}
	}
	m.AzureConfig.VoiceList = voiceList
	m.saveMetadata()
	return nil
}

func (m *UserSettings) ExportMapping() FieldsMapping {
	return FieldsMapping{}
}

func (m *UserSettings) SettingsPathsPortable() bool {

	for _, listPath := range m.WordLists.WordLists {
		if filepath.IsAbs(listPath) {
			return false
		}
	}

	for _, dict := range m.DictionariesConfig.Dicts {
		if filepath.IsAbs(dict.Path) {
			return false
		}
	}
	return true
}

func (m *UserSettings) FixSettingsPaths() error {
	// This should be much more safe

	// For path in dicts fix path
	for name, listPath := range m.WordLists.WordLists {
		if filepath.IsAbs(listPath) {
			relPath, err := filepath.Rel(ConfigDir(), listPath)
			if err != nil {
				return err
			}
			m.SaveList(name, relPath)
		}
	}

	for name, dict := range m.DictionariesConfig.Dicts {
		if filepath.IsAbs(dict.Path) {
			relPath, err := filepath.Rel(ConfigDir(), dict.Path)
			if err != nil {
				return err
			}
			m.SaveDict(name, relPath, dict.Language)
		}
	}

	// For path in userDicts fix path
	return nil
}

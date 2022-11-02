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

type UserSettings struct {
	path        string
	Ran         int
	Dicts       map[string]Dict
	PrimaryDict string
	// We are just saving a few settings here so its fine to
	// Have the values as strings and do conversions
	AutoAdvanceSentence   bool
	PopulateEnglish       bool
	PopulateChinese       bool
	AutoAdvanceEnglish    bool
	AutoAdvanceImage      bool
	GenerateTermAudio     bool
	GenerateSentenceAudio bool
	AutoAdvanceCard       bool

	ShowDefinitions  bool
	EnableChinese    bool
	AzureApiKey      string
	AzureImageApiKey string
	KnownInterval    int

	OnlyFavorites bool
}

func defaultSettings(path string) *UserSettings {
	return &UserSettings{
		path:                  path,
		Ran:                   0,
		Dicts:                 map[string]Dict{},
		PrimaryDict:           "",
		AutoAdvanceSentence:   true,
		PopulateEnglish:       false,
		PopulateChinese:       false,
		AutoAdvanceEnglish:    false,
		AutoAdvanceImage:      false,
		GenerateTermAudio:     false,
		GenerateSentenceAudio: false,
		AutoAdvanceCard:       true,
		ShowDefinitions:       true,
		EnableChinese:         true,
		AzureApiKey:           "",
		AzureImageApiKey:      "",
		KnownInterval:         10,
		OnlyFavorites:         false,
	}
}

var userSettings *UserSettings

func LoadMetadata(path string) (*UserSettings, error) {
	userSettings = defaultSettings(path)

	if _, err := os.Stat(path); err == nil {
		// metadata already exists, read from it
		b, err := os.ReadFile(path)
		if err = json.Unmarshal(b, userSettings); err != nil {
			return nil, err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// metadata does *not* exist, write the default settings
		saveMetadata()
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.Fatal(err)
	}
	return userSettings, nil
}

func saveMetadata() error {
	path := userSettings.path
	str, err := json.MarshalIndent(userSettings, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, str, 0666); err != nil {
		return err
	}
	return nil
}

func UpdateTimesRan() {
	log.Println(userSettings)
	userSettings.Ran += 1
	saveMetadata()
}

func GetTimesRan() int {
	return userSettings.Ran
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
}

func (m *UserSettings) SetUserSettingBool(key string, val bool) {
	setUserSetting(m, key, val)
}

func (m *UserSettings) SetUserSettingInt(key string, val int) {
	setUserSetting(m, key, val)
}

func setUserSetting[T int | string | bool](m *UserSettings, key string, val T) {
	err := reflections.SetField(m, key, val)
	if err != nil {
		log.Fatal(err)
	}
	saveMetadata()
}

func SaveDict(name string, dictPath string, language string) {
	userSettings.Dicts[name] = Dict{
		Path:     dictPath,
		Language: language,
	}
	saveMetadata()
}

func DeleteDict(name string) {
	delete(userSettings.Dicts, name)
	saveMetadata()
}

func LoadDicts() map[string]Dict {
	return userSettings.Dicts
}

func SetPrimaryDict(dictName string) {
	// TODO Make sure its a real dict
	userSettings.PrimaryDict = dictName
	saveMetadata()
}

func GetPrimaryDict() string {
	return userSettings.PrimaryDict
}

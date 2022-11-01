package backend

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"encoding/json"
)

type Dict struct {
	Path     string
	Language string
}

type Metadata struct {
	path        string
	Ran         int
	Dicts       map[string]Dict
	PrimaryDict string
	// We are just saving a few settings here so its fine to
	// Have the values as strings and do conversions
	UserSettings UserSettingsMap
}

type UserSettingsMap struct {
	Settings map[string]UserSetting
}

type TempSettingsMap struct {
	Settings map[string]UserSetting
}

// We setup something to unmarshal the map of settings, to
// prevent the default behaviour of overwriting the entire map
// entry
func (m *UserSettingsMap) UnmarshalJSON(b []byte) error {
	tempMap := &TempSettingsMap{}
	if err := json.Unmarshal(b, tempMap); err != nil {
		return nil
	}
	for name, setting := range tempMap.Settings {
		mSetting, ok := m.Settings[name]
		if !ok {
			return errors.New("Mismatch keys")
		}
		mSetting.Value = setting.Value
		m.Settings[name] = mSetting
	}
	return nil
}

// This is the interface that will be accessed by the frontend
type UserSettings struct {
	// AutoAdvanceSentence   bool
	// PopulateEnglish       bool
	// PopulateChinese       bool
	// AutoAdvanceEnglish    bool
	// AutoAdvanceImage      bool
	// GenerateTermAudio     bool
	// GenerateSentenceAudio bool
	// AutoAdvanceCard       bool
	//
	// ShowDefinitions       bool
	// EnableChinese         bool
	// AzureApiKey           string
	// AzureImageApiKey      string
	// KnownInterval         int
	// OnlyFavorites         bool
}

type UserSetting struct {
	Datatype string `json:"-"` // one of string int bool
	Selector string `json:"-"` // checkbox textbox slider
	Value    string
	Label    string `json:"-"`
	Tooltip  string `json:"-"`
	Group    string `json:"-"`
}

func checkBox(group string, label string, tooltip string, defaultValue bool) UserSetting {
	return UserSetting{
		Datatype: "bool",
		Selector: "checkbox",
		Value:    strconv.FormatBool(defaultValue),
		Label:    label,
		Tooltip:  tooltip,
		Group:    group,
	}
}

func textBox(group string, label string, tooltip string, defaultValue string) UserSetting {
	return UserSetting{
		Datatype: "string",
		Selector: "textbox",
		Value:    defaultValue,
		Label:    label,
		Tooltip:  tooltip,
		Group:    group,
	}
}

func slider(group string, label string, tooltip string, defaultValue int64) UserSetting {
	return UserSetting{
		Datatype: "int",
		Selector: "slider",
		Value:    strconv.FormatInt(defaultValue, 10),
		Label:    label,
		Tooltip:  tooltip,
		Group:    group,
	}
}

func defaultSettings(path string) *Metadata {
	return &Metadata{
		path:        path,
		Ran:         0,
		Dicts:       map[string]Dict{},
		PrimaryDict: "",
		UserSettings: UserSettingsMap{
			Settings: map[string]UserSetting{
				// AutoAdvanceSentence   bool
				"AutoAdvanceSentence": checkBox(
					"CardCreation",
					"Auto advance after sentence selection",
					"After picking a sentence, move to the next step",
					true,
				),
				// PopulateEnglish       bool
				"PopulateEnglish": checkBox(
					"CardCreation",
					"Auto fill english definitions",
					"If only one definition exists, auto select it",
					false,
				),
				// PopulateChinese       bool
				"PopulateChinese": checkBox(
					"CardCreation",
					"Auto fill chinese definitions",
					"If only one definition exists, auto select it",
					false,
				),
				// AutoAdvanceEnglish    bool
				"AutoAdvanceEnglish": checkBox(
					"CardCreation",
					"Auto advance after definition selection",
					"After picking a definition, move to the next step",
					false,
				),
				// AutoAdvanceImage      bool
				"AutoAdvanceImage": checkBox(
					"CardCreation",
					"Auto advance after image selection",
					"After picking a image, move to the next step",
					false,
				),
				// GenerateTermAudio     bool
				"GenerateTermAudio": checkBox(
					"CardCreation",
					"Auto generate audio for keyword",
					"Not implemented yet",
					false,
				),
				// GenerateSentenceAudio bool
				"GenerateSentenceAudio": checkBox(
					"CardCreation",
					"Auto generate audio for example sentence",
					"Not implemented yet",
					false,
				),
				// AutoAdvanceCard       bool
				"AutoAdvanceCard": checkBox(
					"CardCreation",
					"Create card once all fields have been filled",
					"Create card once all fields have been filled",
					true,
				),

				// ShowDefinitions       bool
				"ShowDefinitions": checkBox(
					"Dictionaries",
					"Show Definitions",
					"Show the definitions for words in various tables",
					true,
				),
				// EnableChinese         bool
				"EnableChinese": checkBox(
					"Dictionaries",
					"Use Chinese definitions",
					"Allow flashcards to use chinese definitions instead of just english ones",
					true,
				),
				// AzureApiKey           string
				"AzureApiKey": textBox(
					"Dictionaries",
					"Azure Audio Api Key",
					"Setup an free azure tts account and put your key here",
					"",
				),
				// AzureImageApiKey      string
				"AzureImageApiKey": textBox(
					"Dictionaries",
					"Azure Image Api Key",
					"Setup an free azure bing search and put your key here",
					"",
				),
				// KnownInterval         int
				"KnownInterval": slider(
					"Dictionaries",
					"Time before a word is considered 'known'",
					"How long of an interval in anki before a word is included in generated sentences",
					10,
				),
				// OnlyFavorites         bool
				"OnlyFavorites": checkBox(
					"BookLibrary",
					"Only show favorited books",
					"",
					false,
				),
			},
		},
	}
}

var userSettings *Metadata

func LoadMetadata(path string) error {
	userSettings = defaultSettings(path)

	if _, err := os.Stat(path); err == nil {
		// metadata already exists, read from it
		b, err := os.ReadFile(path)
		if err = json.Unmarshal(b, userSettings); err != nil {
			return err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// metadata does *not* exist, write the default settings
		saveMetadata()
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.Fatal(err)
	}
	return nil
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

func (settings *UserSettings) GetGroupSettings(groupName string) []UserSetting {
	group := []UserSetting{}
	for _, setting := range userSettings.UserSettings.Settings {
		if setting.Group == groupName {
			group = append(group, setting)
		}
	}
	return group
}

func (settings *UserSettings) GetOptionValue(key string) string {
	return GetOptionValue(key)
}

func GetOptionValue(key string) string {
	setting, ok := userSettings.UserSettings.Settings[key]
	if ok {
		return setting.Value
	} else {
		log.Fatal("Bad OptionValue", key)
		return ""
	}
}

func GetOptionValueBool(key string) bool {
	res, _ := strconv.ParseBool(GetOptionValue(key))
	return res
}

func GetOptionValueInt(key string) int {
	res, _ := strconv.ParseInt(GetOptionValue(key), 10, 32)
	return int(res)
}

func SetOptionValue[V int | bool | string](key string, value V) {
	stringVal := fmt.Sprintf("%v", value)
	setting, ok := userSettings.UserSettings.Settings[key]
	if ok {
		setting.Value = stringVal
		userSettings.UserSettings.Settings[key] = setting
		saveMetadata()
	} else {
		log.Println("Bad OptionValue", key)
	}
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

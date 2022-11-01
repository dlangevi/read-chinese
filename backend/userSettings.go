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
	UserSettings map[string]UserSetting
}

type UserSetting struct {
	datatype string // one of string int bool
	selector string // checkbox textbox slider
	Value    string
	label    string
	tooltip  string
	group    string
}

func checkBox(group string, label string, tooltip string, defaultValue bool) UserSetting {
	return UserSetting{
		datatype: "bool",
		selector: "checkbox",
		Value:    strconv.FormatBool(defaultValue),
		label:    label,
		tooltip:  tooltip,
		group:    group,
	}
}

func textBox(group string, label string, tooltip string, defaultValue string) UserSetting {
	return UserSetting{
		datatype: "string",
		selector: "textbox",
		Value:    defaultValue,
		label:    label,
		tooltip:  tooltip,
		group:    group,
	}
}

func slider(group string, label string, tooltip string, defaultValue int64) UserSetting {
	return UserSetting{
		datatype: "int",
		selector: "slider",
		Value:    strconv.FormatInt(defaultValue, 10),
		label:    label,
		tooltip:  tooltip,
		group:    group,
	}
}

func defaultSettings(path string) *Metadata {
	return &Metadata{
		path:        path,
		Ran:         0,
		Dicts:       map[string]Dict{},
		PrimaryDict: "",
		UserSettings: map[string]UserSetting{
			// AutoAdvanceSentence   bool
			"CardCreation": checkBox(
				"AutoAdvanceSentence",
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

// TODO variable
func GetOptionValue(key string) string {
	setting, ok := userSettings.UserSettings[key]
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
	setting, ok := userSettings.UserSettings[key]
	if ok {
		setting.Value = stringVal
		userSettings.UserSettings[key] = setting
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

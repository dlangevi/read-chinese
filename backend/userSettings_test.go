package backend

import (
	"os"
	"path"
	"testing"
)

func TestDefault(t *testing.T) {
	tmpMetaData := path.Join(os.TempDir(), "metadata.json")
	LoadMetadata(tmpMetaData)
	settingName := "EnableChinese"
	setting, ok := userSettings.UserSettings[settingName]
	if !ok {
		t.Errorf("Failed to load basic setting %v", settingName)
	}
	if setting.Group != "Dictionaries" {
		t.Errorf("%v, had the wrong group %#v", settingName, setting)
	}
	if setting.Value != "true" {
		t.Errorf("%v, had the wrong value %v", settingName, setting.Value)
	}
	if userSettings.PrimaryDict != "" {
		t.Errorf("PrimaryDict had wrong default value %v", userSettings.PrimaryDict)
	}
	// TODO have a test where the dictionary does not exist
	SetPrimaryDict("foo")
	if GetPrimaryDict() != "foo" {
		t.Errorf("PrimaryDict was not updated")
	}
	if GetOptionValueBool(settingName) != true {
		t.Errorf("%v was not initialized correctly", settingName)
	}
	SetOptionValue(settingName, false)
	if GetOptionValueBool(settingName) != false {
		t.Errorf("%v was not updated", settingName)
	}
	os.Remove(tmpMetaData)
}

func TestGroupSettings(t *testing.T) {
	tmpMetaData := path.Join(os.TempDir(), "metadatagroup.json")
	LoadMetadata(tmpMetaData)
	settings := &UserSettings{}
	groups := settings.GetGroupSettings("CardCreation")
	if len(groups) != 8 {
		t.Errorf("Groups not good %v", len(groups))
	}
	os.Remove(tmpMetaData)
}

func TestLoadTwice(t *testing.T) {
	tmpMetaData := path.Join(os.TempDir(), "metadatatwice.json")
	LoadMetadata(tmpMetaData)
	settingName := "EnableChinese"
	settingA, _ := userSettings.UserSettings[settingName]
	LoadMetadata(tmpMetaData)
	settingB, _ := userSettings.UserSettings[settingName]
	if settingA != settingB {
		t.Errorf("Load twice %v, %v", settingA, settingB)
	}
	os.Remove(tmpMetaData)
}

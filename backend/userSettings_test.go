package backend

import (
	"os"
	"path"
	"testing"
)

func TestDefault(t *testing.T) {
	tmpMetaData := path.Join(os.TempDir(), "metadata.json")
	userSettings, _ := LoadMetadata(tmpMetaData)

	settingName := "EnableChinese"
	setting := userSettings.EnableChinese
	if setting != true {
		t.Errorf("%v, had the wrong value %v", settingName, setting)
	}
	if userSettings.PrimaryDict != "" {
		t.Errorf("PrimaryDict had wrong default value %v", userSettings.PrimaryDict)
	}
	// TODO have a test where the dictionary does not exist
	userSettings.SetPrimaryDict("foo")
	if userSettings.PrimaryDict != "foo" {
		t.Errorf("PrimaryDict was not updated")
	}
	if userSettings.GetUserSettingBool(settingName) != true {
		t.Errorf("%v was not initialized correctly", settingName)
	}
	userSettings.SetUserSettingBool(settingName, false)
	if userSettings.GetUserSettingBool(settingName) != false {
		t.Errorf("%v was not updated", settingName)
	}
	os.Remove(tmpMetaData)
}

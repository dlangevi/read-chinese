package backend

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	tmpMetaData := path.Join(os.TempDir(), "metadata.json")
	userSettings, _ := LoadMetadata(tmpMetaData, testRuntime)

	settingName := "EnableChinese"
	setting := userSettings.DictionariesConfig.EnableChinese
	if setting != false {
		t.Errorf("%v, had the wrong value %v", settingName, setting)
	}
	if userSettings.DictionariesConfig.PrimaryDict != "" {
		t.Errorf("PrimaryDict had wrong default value %v",
			userSettings.DictionariesConfig.PrimaryDict)
	}
	// TODO have a test where the dictionary does not exist
	userSettings.SetPrimaryDict("foo")
	if userSettings.DictionariesConfig.PrimaryDict != "foo" {
		t.Errorf("PrimaryDict was not updated")
	}
	os.Remove(tmpMetaData)
}

func TestNew(t *testing.T) {
	newConfig := defaultConfig("foobar")
	assert.Equal(t, true, getValue(newConfig, "AddProgramTag"))
	assert.Equal(t, false, getValue(newConfig, "OnlyFavorites"))
	err := setValue(newConfig, "OnlyFavorites", true)
	assert.Nil(t, err)
	assert.Equal(t, true, getValue(newConfig, "OnlyFavorites"))

}

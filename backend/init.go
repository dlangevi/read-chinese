package backend

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type Backend struct {
	RuntimeContext *context.Context

	// Primary data layer
	UserSettings *UserConfig
	DB           *sqlx.DB

	// Independent Libraries
	ImageClient *ImageClient

	// Libraries required by other Libraries
	BookLibrary   BookLibrary
	KnownWords    *KnownWords
	Dictionaries  *Dictionaries
	AnkiInterface AnkiInterface

	// Libraries that require other Libraires
	Segmentation *Segmentation
	Generator    *Generator
	Calibre      *Calibre
}

func (b *Backend) HealthCheck() (string, error) {
	checkBooks, err := b.BookLibrary.HealthCheck()
	if err != nil {
		return "", err
	}
	if checkBooks != "" {
		return checkBooks, nil
	}
	checkDicts := b.Dictionaries.HealthCheck()
	if checkDicts != "" {
		return checkDicts, nil
	}

	checkAnki := b.AnkiInterface.HealthCheck()
	if checkAnki != "" {
		return checkAnki, nil
	}

	configureAnki, err := b.AnkiInterface.ConfigurationCheck()
	if err != nil {
		return "", err
	}
	if configureAnki != "" {
		return configureAnki, nil
	}

	return "", nil
}

func StartBackend(ctx *context.Context,
	sqlPath string,
	metadataPath string) (*Backend, error) {

	db, err := NewDB(sqlPath)
	if err != nil {
		return nil, err
	}
	err = RunMigrateScripts(db)
	if err != nil {
		return nil, err
	}
	userSettings, err := LoadMetadata(metadataPath)
	if err != nil {
		return nil, err
	}

	userSettings.UpdateTimesRan()
	ran := userSettings.GetTimesRan()
	log.Printf("Ran %v times", ran)

	runtime := &Backend{
		RuntimeContext: ctx,
		UserSettings:   userSettings,
		DB:             db,

		KnownWords:  NewKnownWords(db, userSettings),
		ImageClient: NewImageClient(userSettings),
	}

	runtime.Dictionaries = NewDictionaries(userSettings, runtime.KnownWords)
	runtime.AnkiInterface = NewAnkiInterface(userSettings, runtime.KnownWords)

	err = UnloadJiebaDicts()
	if err != nil {
		return nil, err
	}
	s, err := NewSegmentation(runtime.Dictionaries)
	if err != nil {
		return nil, err
	}

	runtime.BookLibrary = NewBookLibrary(db, s, runtime.KnownWords)
	runtime.Segmentation = s
	runtime.Generator = NewGenerator(userSettings, s, runtime.BookLibrary, runtime.KnownWords)
	runtime.Calibre = NewCalibre(runtime.BookLibrary, runtime.Generator)

	return runtime, nil
}

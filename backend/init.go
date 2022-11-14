package backend

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type Backend struct {
	RuntimeContext *context.Context

	// Primary data layer
	UserSettings *UserSettings
	DB           *sqlx.DB

	// Independent Libraries
	ImageClient *ImageClient

	// Libraries required by other Libraries
	BookLibrary   *BookLibrary
	KnownWords    *KnownWords
	Dictionaries  *Dictionaries
	AnkiInterface *AnkiInterface

	// Libraries that require other Libraires
	Segmentation *Segmentation
	Generator    *Generator
	Calibre      *Calibre
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

	s, err := NewSegmentation(runtime.Dictionaries)
	if err != nil {
		return nil, err
	}

	runtime.BookLibrary = NewBookLibrary(db, s, runtime.KnownWords)
	runtime.Segmentation = s
	runtime.Generator = NewGenerator(s, runtime.BookLibrary, runtime.KnownWords)
	runtime.Calibre = NewCalibre(runtime.BookLibrary)

	return runtime, nil
}

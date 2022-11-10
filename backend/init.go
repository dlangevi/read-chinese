package backend

import (
	"context"
	"log"
)

type Backend struct {
	RuntimeContext *context.Context
	BookLibrary    *BookLibrary
	KnownWords     *KnownWords
	UserSettings   *UserSettings
	ImageClient    *ImageClient
	Dictionaries   *Dictionaries
	Segmentation   *Segmentation
	Generator      *Generator
	AnkiInterface  *AnkiInterface
	Calibre        *Calibre
}

var runtime *Backend

func StartBackend(ctx *context.Context) *Backend {

	err := NewDB(ConfigDir("db.sqlite3"))
	if err != nil {
		log.Fatal(err)
	}
	err = RunMigrateScripts()
	if err != nil {
		log.Fatal(err)
	}
	userSettings, err = LoadMetadata(ConfigDir("newmetadata.json"))
	if err != nil {
		log.Fatal(err)
	}

	UpdateTimesRan()
	ran := GetTimesRan()
	log.Printf("Ran %v times", ran)

	d := NewDictionaries()
	s, err := NewSegmentation(d)
	if err != nil {
		log.Fatal(err)
	}

	runtime = &Backend{
		RuntimeContext: ctx,
		BookLibrary:    &BookLibrary{},
		KnownWords:     NewKnownWords(),
		UserSettings:   userSettings,
		ImageClient:    &ImageClient{},
		Dictionaries:   d,
		Segmentation:   s,
		Generator:      &Generator{},
		AnkiInterface:  NewAnkiInterface(),
		Calibre:        &Calibre{},
	}

	return runtime
}

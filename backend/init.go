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
	Segmentation   *Segmentation
	ImageClient    *ImageClient
	Dictionaries   *Dictionaries
	Generator      *Generator
	AnkiInterface  *AnkiInterface
	Calibre        *Calibre
}

func StartBackend(ctx *context.Context,
	sqlPath string,
	metadataPath string) (*Backend, error) {

	err := NewDB(sqlPath)
	if err != nil {
		return nil, err
	}
	err = RunMigrateScripts()
	if err != nil {
		return nil, err
	}
	userSettings, err = LoadMetadata(metadataPath)
	if err != nil {
		return nil, err
	}

	UpdateTimesRan()
	ran := GetTimesRan()
	log.Printf("Ran %v times", ran)

	d := NewDictionaries()
	s, err := NewSegmentation(d)
	if err != nil {
		return nil, err
	}

	b := NewBookLibrary(s)

	runtime := &Backend{
		RuntimeContext: ctx,
		BookLibrary:    b,
		KnownWords:     NewKnownWords(),
		UserSettings:   userSettings,
		Segmentation:   s,
		ImageClient:    NewImageClient(),
		Dictionaries:   d,
		Generator:      NewGenerator(s),
		AnkiInterface:  NewAnkiInterface(),
		Calibre:        NewCalibre(b),
	}

	return runtime, nil
}

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
	TextToSpeech   *TextToSpeech
	Dictionaries   *Dictionaries
	Segmentation   *Segmentation
	Generator      *Generator
	AnkiInterface  *AnkiInterface
	Calibre        *Calibre
}

var runtime *Backend

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

	runtime = &Backend{
		RuntimeContext: ctx,
		BookLibrary:    &BookLibrary{},
		KnownWords:     NewKnownWords(),
		UserSettings:   userSettings,
		TextToSpeech:   NewTextToSpeach(),
		ImageClient:    NewImageClient(),
		Dictionaries:   d,
		Segmentation:   s,
		Generator:      &Generator{},
		AnkiInterface:  NewAnkiInterface(),
		Calibre:        &Calibre{},
	}

	return runtime, nil
}

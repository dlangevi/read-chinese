package backend

import (
	"context"
	"fmt"
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
}

var runtime *Backend

func StartBackend(ctx *context.Context) *Backend {

	err := NewDB("/home/dlangevi/.config/read-chinese/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	err = RunMigrateScripts()
	if err != nil {
		log.Fatal(err)
	}
	userSettings, err = LoadMetadata("/home/dlangevi/.config/read-chinese/newmetadata.json")
	if err != nil {
		log.Fatal(err)
	}

	UpdateTimesRan()
	ran := GetTimesRan()
	fmt.Println("Ran {} times", ran)

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
	}

	return runtime
}

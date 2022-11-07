package backend

import (
	"context"
	"fmt"
	"log"
	"read-chinese/backend/segmentation"
)

type Backend struct {
	RuntimeContext *context.Context
	BookLibrary    *BookLibrary
	KnownWords     *KnownWords
	UserSettings   *UserSettings
	ImageClient    *ImageClient
	Dictionaries   *Dictionaries
	Segmentation   *segmentation.Segmentation
	Generator      *Generator
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

	s, err := segmentation.NewSegmentation()
	if err != nil {
		log.Fatal(err)
	}

	runtime = &Backend{
		RuntimeContext: ctx,
		BookLibrary:    &BookLibrary{},
		KnownWords:     NewKnownWords(),
		UserSettings:   userSettings,
		ImageClient:    &ImageClient{},
		Dictionaries:   NewDictionaries(),
		Segmentation:   s,
		Generator:      &Generator{},
	}

	return runtime
}

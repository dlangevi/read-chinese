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
}

func StartBackend(ctx *context.Context) *Backend {

	err := NewDB("/home/dlangevi/.config/read-chinese/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	err = RunMigrateScripts()
	if err != nil {
		log.Fatal(err)
	}
	err = LoadMetadata("/home/dlangevi/.config/read-chinese/newmetadata.json")
	if err != nil {
		log.Fatal(err)
	}
	UpdateTimesRan()
	ran := GetTimesRan()
	fmt.Println("Ran {} times", ran)

	// settings, err := LoadSettings()
	// if err != nil {
	// 	log.Error().Msg("Failed to load settings")
	// }
	return &Backend{
		RuntimeContext: ctx,
		BookLibrary:    &BookLibrary{},
		KnownWords:     NewKnownWords(),
	}
}

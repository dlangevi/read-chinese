package core

import (
  "log"

	"github.com/adrg/xdg"
)

func configDir() string {
	configDirPath, err := xdg.ConfigFile("read-chinese/")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config Dir at:", configDirPath)
	return configDirPath
}

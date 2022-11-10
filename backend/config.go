package backend

import (
	"log"
	"path"

	"github.com/adrg/xdg"
)

func ConfigDir(file ...string) string {
	file = append([]string{"read-chinese"}, file...)
	local := path.Join(file...)
	configDirPath, err := xdg.ConfigFile(local)
	if err != nil {
		log.Fatal(err)
	}
	return configDirPath
}

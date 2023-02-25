package backend

import (
	"fmt"
	"log"
	"path"

	"github.com/adrg/xdg"
)

var configFolderName = "read-chinese"

func SetTestUser(username string) {
	configFolderName = fmt.Sprintf("test-read-chinese-%v", username)
}

func ConfigDir(file ...string) string {
	file = append([]string{configFolderName}, file...)
	local := path.Join(file...)
	configDirPath, err := xdg.ConfigFile(local)
	if err != nil {
		log.Fatal(err)
	}
	return configDirPath
}

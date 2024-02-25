package backend

import (
	"fmt"
	"log"
	"path"

	"github.com/adrg/xdg"

	"archive/zip"
	"io"
	"os"
	"path/filepath"
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

type ConfigExposer struct {
}

func (c *ConfigExposer) GetConfigDir() string {
	return ConfigDir()
}

// Zips "./input" into "./output.zip"
func (c *ConfigExposer) ZipConfigDir() {
	// TODO this does not work with the unicode paths I have been using.
	// Maybe v2 would be to reconvert everthing to using ascii (via pinyin)
	file, err := os.Create(ConfigDir("archive.zip"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		relpath, err := filepath.Rel(ConfigDir(), path)
		if err != nil {
			return err
		}

		h := &zip.FileHeader{Name: relpath, Method: zip.Deflate, Flags: 0x800}
		f, err := w.CreateHeader(h)
		// f, err := w.Create(relpath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	// This needs to be updated whenever we add a new storage location
	var configSubDirs = []string{
		"bookRawText",
		"userDicts",
		"covers",
		"jiebaDicts",
		"segmentationCache",
		"userLists",
	}

	for _, subpath := range configSubDirs {
		err = filepath.Walk(ConfigDir(subpath), walker)
		if err != nil {
			panic(err)
		}
	}
	// TODO handle db.sqlite3-wal and db.sqlite3-shm
	// db.sqlite3   metadata.json
}

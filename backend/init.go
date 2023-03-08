package backend

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Backend struct {
	sqlPath      string
	metadataPath string

	ctx context.Context

	// Primary data layer
	UserSettings *UserSettings
	DB           *sqlx.DB

	// Independent Libraries
	ImageClient *ImageClient

	// Libraries required by other Libraries
	BookLibrary   BookLibrary
	KnownWords    *KnownWords
	Dictionaries  *Dictionaries
	AnkiInterface AnkiInterface

	// Libraries that require other Libraires
	Segmentation *Segmentation
	Generator    *Generator
	Calibre      *Calibre
}

func (b *Backend) HealthCheck() error {
	err := b.BookLibrary.HealthCheck()
	if err != nil {
		return err
	}
	err = b.Dictionaries.HealthCheck()
	if err != nil {
		return err
	}

	err = b.AnkiInterface.HealthCheck()
	if err != nil {
		return err
	}

	err = b.AnkiInterface.ConfigurationCheck()
	if err != nil {
		return err
	}
	return nil
}

func NewBackend(
	sqlPath string,
	metadataPath string) *Backend {

	db, err := NewDB(sqlPath)
	if err != nil {
		log.Fatal(err)
	}
	err = RunMigrateScripts(db)
	if err != nil {
		log.Fatal(err)
	}

	// The order of setup may be important
	backend := &Backend{
		sqlPath:      sqlPath,
		metadataPath: metadataPath,
		DB:           db,
	}

	userSettings, err := LoadMetadata(metadataPath, backend)
	if err != nil {
		log.Fatal(err)
	}

	userSettings.UpdateTimesRan()
	ran := userSettings.GetTimesRan()
	log.Printf("Ran %v times", ran)

	backend.UserSettings = userSettings
	backend.KnownWords = NewKnownWords(db, userSettings)
	backend.ImageClient = NewImageClient(userSettings)
	backend.Dictionaries = NewDictionaries(backend)
	backend.AnkiInterface = NewAnkiInterface(backend)

	err = UnloadJiebaDicts()
	if err != nil {
		log.Fatal(err)
	}
	s, err := NewSegmentation(backend.Dictionaries)
	if err != nil {
		log.Fatal(err)
	}

	backend.BookLibrary = NewBookLibrary(backend, db)
	backend.Segmentation = s
	backend.Generator = NewGenerator(backend)
	backend.Calibre = NewCalibre(backend)

	return backend
}

func (b *Backend) Startup(ctx context.Context) {
	b.ctx = ctx

}

func (b *Backend) SaveFile() (string, error) {
	selectedFile, err := runtime.SaveFileDialog(b.ctx, runtime.SaveDialogOptions{
		Title: "Save book stats",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "csv file",
				Pattern:     "*.csv",
			},
		},
		DefaultFilename: "BookStats.csv",
	})
	return selectedFile, err
}
func (b *Backend) FilePicker(extension string) (string, error) {
	log.Println("requesting file")
	selectedFile, err := runtime.OpenFileDialog(b.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: extension,
				Pattern:     fmt.Sprintf("*.%v", extension),
			},
		},
	})
	return selectedFile, err
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (b *Backend) DomReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (b *Backend) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (b *Backend) Shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}

type ReportedDownload struct {
	backend          *Backend
	TotalBytes       uint64
	lastSentProgress uint64
	bytesSoFar       uint64
}

func (wc *ReportedDownload) Write(p []byte) (int, error) {
	n := len(p)
	wc.bytesSoFar += uint64(n)
	progress := wc.bytesSoFar * 100 / wc.TotalBytes
	if progress > wc.lastSentProgress {
		wc.backend.progress()
		wc.lastSentProgress = progress
	}
	return n, nil
}

func (b *Backend) setupProgress(
	message string,
	steps int) {
	runtime.EventsEmit(b.ctx, "setupProgress", message, steps)
}

func (b *Backend) progress() {
	runtime.EventsEmit(b.ctx, "progress")
}

package main

// #cgo LDFLAGS: -static-libstdc++ -static-libgcc
import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"read-chinese/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

// var icon []byte

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedFilename := req.URL.Path
	println("Requesting file:", requestedFilename)
	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	logFile := flag.String("log", "", "log file")
	flag.Parse()

	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}
	// Create an instance of the app structure
	// 创建一个App结构体实例
	log.SetFlags(log.Ltime | log.Lshortfile)

	app := NewApp()
	sqlDbPath := backend.ConfigDir("db.sqlite3")
	metadataPath := backend.ConfigDir("metadata.json")
	backend, err := backend.StartBackend(&app.ctx, sqlDbPath, metadataPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create application with options
	// 使用选项创建应用
	err = wails.Run(&options.App{
		Title:             "read-chinese",
		Width:             1200,
		Height:            800,
		MinWidth:          1200,
		MinHeight:         800,
		MaxWidth:          1920,
		MaxHeight:         1080,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.ERROR,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnBeforeClose:     app.beforeClose,
		OnShutdown:        app.shutdown,
		WindowStartState:  options.Maximised,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    NewFileLoader(),
			Middleware: nil,
		},
		Bind: []interface{}{
			app,
			backend.BookLibrary,
			backend.KnownWords,
			backend.UserSettings,
			backend.ImageClient,
			backend.Dictionaries,
			backend.Generator,
			backend.AnkiInterface,
			backend.Calibre,
		},
		// Windows platform specific options
		// Windows平台特定选项
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
			// User messages that can be customised
			// Messages: &windows.Messages{
			// 	InstallationRequired: "",
			// 	UpdateRequired:       "",
			// 	MissingRequirements:  "",
			// 	Webview2NotInstalled: "",
			// 	Error:                "",
			// 	FailedToInstall:      "",
			// 	DownloadPage:         "",
			// 	PressOKToInstall:     "",
			// 	ContactAdmin:         "",
			// 	InvalidFixedWebview2: "",
			// },
		},
		// Mac platform specific options
		// Mac平台特定选项
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
				HideToolbarSeparator:       false,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "read-chinese",
				Message: "Read more learn more read more",
				// Icon:    icon,
			},
		},
		Linux: &linux.Options{
			// Icon: icon,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

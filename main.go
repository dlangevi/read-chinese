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
	coverLocation := backend.ConfigDir("covers", requestedFilename)
	fileData, err := os.ReadFile(coverLocation)
	if err != nil {
		log.Println("Error reading file", err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	logFile := flag.String("log", "", "log file")
	testUser := flag.String(
		"user", "",
		"-user {user} will run the application with a different profile")
	flag.Parse()

	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	if *testUser != "" {
		log.Println("Running as test user: ", *testUser)
		backend.SetTestUser(*testUser)
	}
	// Create an instance of the app structure
	// 创建一个App结构体实例
	log.SetFlags(log.Ltime | log.Lshortfile)

	sqlDbPath := backend.ConfigDir("db.sqlite3")
	metadataPath := backend.ConfigDir("metadata.json")
	backendObj := backend.NewBackend(sqlDbPath, metadataPath)

	// Create application with options
	// 使用选项创建应用
	err := wails.Run(&options.App{
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
		OnStartup:         backendObj.Startup,
		OnDomReady:        backendObj.DomReady,
		OnBeforeClose:     backendObj.BeforeClose,
		OnShutdown:        backendObj.Shutdown,
		WindowStartState:  options.Maximised,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    NewFileLoader(),
			Middleware: nil,
		},
		Bind: []interface{}{
			backendObj,
			backendObj.BookLibrary,
			backendObj.KnownWords,
			backendObj.UserSettings,
			backendObj.ImageClient,
			backendObj.TextToSpeech,
			backendObj.AiGenerator,
			backendObj.Dictionaries,
			backendObj.Generator,
			backendObj.AnkiInterface,
			backendObj.Calibre,
			backendObj.WordLists,
			&backend.ConfigExposer{},
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
			WebviewGpuPolicy: linux.WebviewGpuPolicyNever,
			// Icon: icon,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

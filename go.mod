module read-chinese

// module github.com/dlangevi/read-chinese-wails

go 1.18

require (
	github.com/adrg/xdg v0.4.0
	github.com/atselvan/ankiconnect v1.0.1
	github.com/go-resty/resty/v2 v2.7.0
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/jmoiron/sqlx v1.3.5
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/oleiade/reflections v1.0.1
	github.com/privatesquare/bkst-go-utils v1.5.4
	github.com/wailsapp/wails/v2 v2.2.0
	github.com/yanyiwu/gojieba v1.2.0
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
)

require (
	github.com/bep/debounce v1.2.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.7.7 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jarcoal/httpmock v1.0.8 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/labstack/echo/v4 v4.9.0 // indirect
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/leaanthony/go-ansi-parser v1.0.1 // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.5.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/samber/lo v1.27.1 // indirect
	github.com/tkrajina/go-reflector v0.5.5 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// Clone dlangevi/ankiconnect
// To update uncomment and then go mod tidy to download latest
// replace github.com/atselvan/ankiconnect => github.com/dlangevi/ankiconnect main
replace github.com/atselvan/ankiconnect => github.com/dlangevi/ankiconnect v0.0.0-20230127213958-2bd0a1b4c344

// replace github.com/atselvan/ankiconnect => ../ankiconnect

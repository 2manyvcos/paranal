module github.com/2manyvcos/paranal

go 1.27.4

ignore (
	./api/reference
	./api/schema
	./client/dist
	./client/public
	./client/src
	./docs
	node_modules
)

require (
	github.com/adrg/xdg v0.5.3
	github.com/bep/debounce v1.2.1
	github.com/fsnotify/fsnotify v1.10.1
	github.com/go-co-op/gocron/v2 v2.23.1
	github.com/go-viper/mapstructure/v2 v2.5.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/jplorg/jpl/go/v2 v2.1.0
	github.com/spf13/cobra v1.10.2
	go.yaml.in/yaml/v4 v4.0.0
	golang.org/x/crypto v0.59.0
	golang.org/x/mod v0.43.0
	modernc.org/sqlite v1.62.0
)

require (
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.1.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/spf13/pflag v1.1.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	modernc.org/libc v1.81.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
)

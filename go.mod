module github.com/2manyvcos/paranal

go 1.25.0

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
	github.com/go-co-op/gocron/v2 v2.23.0
	github.com/go-viper/mapstructure/v2 v2.5.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/jplorg/jpl/go/v2 v2.1.0
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
)

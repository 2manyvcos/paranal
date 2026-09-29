module github.com/2manyvcos/paranal

go 1.24.4

ignore (
	./api/reference
	./api/schema
	./client/dist
	./client/public
	./client/src
	./docs
	node_modules
)

require github.com/spf13/cobra v1.10.2

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
)

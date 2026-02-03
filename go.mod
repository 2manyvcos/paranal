module github.com/2manyvcos/paranal

go 1.25.3

ignore (
	./client/dist
	./client/public
	./client/src
	./docs
	node_modules
)

require (
	github.com/alexedwards/argon2id v1.0.0
	github.com/go-viper/mapstructure/v2 v2.5.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/uuid v1.6.0
	github.com/jplorg/jpl/go/v2 v2.0.1
	github.com/spf13/cobra v1.10.2
	golang.org/x/crypto v0.47.0
	modernc.org/sqlite v1.44.3
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/exp v0.0.0-20260112195511-716be5621a96 // indirect
	golang.org/x/sys v0.40.0 // indirect
	modernc.org/libc v1.67.7 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

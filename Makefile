include .env

install:
	echo " > Downloading go dependencies"
	go mod download

run:
	echo " > Running proxy checker"
	go run cmd/server/main.go -proxy_url$(PROXY_URL) -timeout=$(PROXY_TIMEOUT)

compile:
	echo " > Proxy checker compiling started [Windows, Linux and MacOS]"
	GOOS=windows GOARCH=amd64 go build -o cmd/server main.go
	GOOS=windows GOARCH=arm64 go build -o cmd/server main.go
	GOOS=linux GOARCH=amd64 go build -o cmd/server main.go
	GOOS=linux GOARCH=arm64 go build -o cmd/server main.go
	GOOS=darwin GOARCH=amd64 go build -o cmd/server main.go
	GOOS=darwin GOARCH=arm64 go build -o cmd/server main.go

all:
	install
	build
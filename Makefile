include .env

help:
	@fgrep -h "#include" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

install: #include
	echo " > Downloading go dependencies"
	go mod download

run: #include
	echo " > Running proxy checker"
	go run cmd/server/main.go

compile: #include
	echo " > Proxy checker compiling started [Windows, Linux and MacOS]"
	GOOS=windows GOARCH=amd64 go build -o cmd/server
	GOOS=windows GOARCH=arm64 go build -o cmd/server
	GOOS=linux GOARCH=amd64 go build -o cmd/server
	GOOS=linux GOARCH=arm64 go build -o cmd/server
	GOOS=darwin GOARCH=amd64 go build -o cmd/server
	GOOS=darwin GOARCH=arm64 go build -o cmd/server

all: #include
	install
	build
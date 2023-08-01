include .env

.PHONY: help
help:
	@fgrep -h "#include" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: install
install: #include
	@echo " > Downloading go dependencies"
	@go mod download
	@go mod tidy
	@echo " > All dependencies successfully installed"

.PHONY: run
run: install #include
	@echo " > Running proxy checker"
	@go run ./cmd/server

.PHONY: compile
compile: install #include
	@echo " > Proxy checker compiling started [Windows, Linux and MacOS]"

	@echo " > Compiling for windows amd64"
	@GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64 ./cmd/server

	@echo " > Compiling for windows arm64"
	@GOOS=windows GOARCH=arm64 go build -o ./bin/windows_arm64 ./cmd/server

	@echo " > Compiling for linux amd64"
	@GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64 ./cmd/server

	@echo " > Compiling for linux arm64"
	@GOOS=linux GOARCH=arm64 go build -o ./bin/linux_arm64 ./cmd/server

	@echo " > Compiling for macos amd64"
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/macos_amd64 ./cmd/server

	@echo " > Compiling for macos arm64"
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/macos_arm64 ./cmd/server
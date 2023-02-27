NAME := $(shell basename "$(PWD)")

.PHONY: build

build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)_linux -ldflags "-X main.Version=$(VERSION)"

build_mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)_mac -ldflags "-X main.Version=$(VERSION)"

build_mac_arm:
	GOOS=darwin GOARCH=arm64 go build -o bin/$(NAME)_mac_arm -ldflags "-X main.Version=$(VERSION)"

build_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(NAME).exe -ldflags "-X main.Version=$(VERSION)"

my_build:
	go build -o bin/$(NAME) -ldflags "-X main.Version=$(VERSION)"

build: build_linux build_mac build_mac_arm build_windows

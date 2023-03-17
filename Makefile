NAME := $(shell basename "$(PWD)")

.PHONY: build

build_linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(NAME)_linux -ldflags "-X main.Version=$(VERSION)"

build_mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(NAME)_mac -ldflags "-X main.Version=$(VERSION)"

build_mac_arm:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/$(NAME)_mac_arm -ldflags "-X main.Version=$(VERSION)"

build_windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(NAME).exe -ldflags "-X main.Version=$(VERSION)"

build:
	go build -o bin/$(NAME) -ldflags "-X main.Version=$(VERSION)"

build_image: build_linux
	docker build -t randome:latest .

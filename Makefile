NAME := $(shell basename "$(PWD)")
IMAGE_VERSION := "0.0.1"

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
	docker build -t arjunmahishi/randome:$(IMAGE_VERSION) .
	docker build -t arjunmahishi/randome:latest .

release: build_image
	docker push arjunmahishi/randome:$(IMAGE_VERSION)
	docker push arjunmahishi/randome:latest

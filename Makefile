GO := go
NAME := quick

VERSION ?= dev

all: build

build:
	$(GO) build \
		-ldflags "-X 'main.BUILD_TIME=$$(date +%Y-%m-%dT%H:%M:%S)' -X 'main.VERSION=$(VERSION)'" \
		-o bin/$(NAME) ./cmd

clean:
	rm -rf bin/

.PHONY: all build clean

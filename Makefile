GO := go
NAME := quick

COVER_OUT := coverage.out
COVER_HTML := coverage.html

VERSION ?= dev

all: build

build:
	$(GO) build \
		-ldflags "-X 'main.BUILD_TIME=$$(date +%Y-%m-%dT%H:%M:%S)' -X 'main.VERSION=$(VERSION)'" \
		-o bin/$(NAME) ./cmd

cover-integration:
	$(GO) test -race -coverprofile=$(COVER_OUT) -coverpkg=./internal/... ./internal/...
	$(GO) tool cover -func=$(COVER_OUT)
	$(GO) tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@echo ""
	@echo "HTML report written to $(COVER_HTML)"

clean:
	rm -rf bin/

.PHONY: all build clean

GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=quote-telegram-bot

all: linux_amd64

linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)
darwin_amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)

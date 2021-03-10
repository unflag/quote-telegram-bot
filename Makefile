GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=quote-telegram-bot
CGO_ENABLED=0

all: linux_amd64

linux_amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)

GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=quote-telegram-bot

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME)

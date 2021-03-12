GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=quote-telegram-bot
VERSION=`git tag`
DATE=`date +%FT%T%z`
GOARCH="amd64"
LDFLAGS="-X main.Name=${BINARY_NAME} -X main.Version=${VERSION:-pipeline.git.tag} -X main.Date=${DATE}"

all: linux_amd64

linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} $(GOBUILD) -ldflags ${LDFLAGS} -o $(BINARY_NAME)

darwin_amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} $(GOBUILD) -ldflags ${LDFLAGS} -o $(BINARY_NAME)

clean:
	if [ -f ${BINARY_NAME} ]; then rm ${BINARY_NAME}; fi
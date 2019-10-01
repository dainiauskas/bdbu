VERSION := $(shell git describe --abbrev=0 --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
NAME := $(shell echo $(PROJECTNAME) | tr '[:lower:]' '[:upper:]')

GOBASE 	:= $(shell pwd)
GOBIN 	:= $(GOBASE)/bin

BINARY := $(GOBIN)/${PROJECTNAME}

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-s -X betsy/cmd.Version=$(VERSION) -X betsy/cmd.Build=$(BUILD) -X betsy/cmd.Name=$(NAME)"

all: before linux win386 win64 after

before:
	cp config.example.yaml $(GOBIN)/config.yaml

linux:
	GOOS=linux go build -tags release $(LDFLAGS) -o $(BINARY)
	upx-ucl $(BINARY)
	zip -j $(BINARY)_$(VERSION)_linux.zip $(BINARY) $(GOBIN)/config.yaml README.md

win64:
	GOOS=windows go build -tags release $(LDFLAGS) -o $(BINARY).exe
	upx-ucl $(BINARY).exe
	zip -j $(BINARY)_$(VERSION)_win.zip $(BINARY).exe $(GOBIN)/config.yaml README.md

win386:
	GOOS=windows GOARCH=386 go build -tags release $(LDFLAGS) -o $(BINARY).exe
	upx-ucl $(BINARY).exe
	zip -j $(BINARY)_$(VERSION)_win_386.zip $(BINARY).exe $(GOBIN)/config.yaml README.md

after:
	@-rm $(GOBIN)/config.yaml 2> /dev/null

clean:
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@-rm $(GOBIN)/$(PROJECTNAME).exe 2> /dev/null

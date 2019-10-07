VERSION := $(shell git describe --abbrev=0 --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
NAME := $(shell echo $(PROJECTNAME) | tr '[:lower:]' '[:upper:]')

GOBASE 	:= $(shell pwd)
GOBIN 	:= $(GOBASE)/bin

BINARY := $(GOBIN)/${PROJECTNAME}

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-s -X $(PROJECTNAME)/cmd.Version=$(VERSION) -X $(PROJECTNAME)/cmd.Build=$(BUILD) -X $(PROJECTNAME)/cmd.Name=$(NAME)"

all: before linux win386 win64 clean source

before:
	cp $(PROJECTNAME).example.yaml $(GOBIN)/$(PROJECTNAME).yaml

linux:
	GOOS=linux go build -tags release $(LDFLAGS) -o $(BINARY)
	upx-ucl $(BINARY)
	zip -j $(BINARY)_$(VERSION)_linux.zip $(BINARY) $(GOBIN)/$(PROJECTNAME).yaml README.md

win64:
	GOOS=windows go build -tags release $(LDFLAGS) -o $(BINARY).exe
	upx-ucl $(BINARY).exe
	zip -j $(BINARY)_$(VERSION)_win.zip $(BINARY).exe $(GOBIN)/$(PROJECTNAME).yaml README.md

win386:
	GOOS=windows GOARCH=386 go build -tags release $(LDFLAGS) -o $(BINARY).exe
	upx-ucl $(BINARY).exe
	zip -j $(BINARY)_$(VERSION)_win_386.zip $(BINARY).exe $(GOBIN)/$(PROJECTNAME).yaml README.md

source:
	zip $(BINARY)_$(VERSION)_src.zip app/*.go cmd/*.go models/*.go *.example.yaml makefile .gitignore

clean:
	@-rm $(GOBIN)/$(PROJECTNAME).yaml 2> /dev/null
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@-rm $(GOBIN)/$(PROJECTNAME).exe 2> /dev/null

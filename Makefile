VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECT := $(shell basename "$(PWD)")

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOOS := "linux"
GOARCH := "amd64"

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(Build)"

install:
	@go get -u

build:
	@echo ">  Building binary"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(PROJECT) *.go

.PHONY: install build
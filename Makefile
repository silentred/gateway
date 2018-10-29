# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

CURDIR := $(shell pwd)
GO        := go
GOBUILD   := $(GO) build
GOTEST    := $(GO) test

OS        := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
PACKAGES  := $$(go list ./...| grep -vE 'vendor|tests')
FILES     := $$(find . -name '*.go' | grep -vE 'vendor')
TARGET	  := "gateway"
LDFLAGS   += -X "main.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS   += -X "main.GitHash=$(shell git rev-parse HEAD)"

test:
	$(GOTEST) $(PACKAGES) -cover

build: build-web build-go

build-web:
	npm --prefix=webui run build
	go-bindata -prefix=webui/dist -o=admin/bindata.go webui/dist/...
	sed -i '' -e 's/package main/package admin/g' admin/bindata.go

build-go:
	$(GOBUILD) -ldflags '$(LDFLAGS)' -o $(TARGET)

dev: test build

clean:
	rm $(TARGET)

devui:
	npm --prefix=webui run dev

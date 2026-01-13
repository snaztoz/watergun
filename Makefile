VERSION 	?= dev
GIT_COMMIT 	:= $(shell git rev-parse --short HEAD)
BUILD_DATE 	?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

GOOS   		?= $(shell go env "GOOS")
GOARCH 		?= $(shell go env "GOARCH")

DIST_DIR 	:= dist
TMP_DIR		:= tmp
MAIN_FILE 	:= watergun.go

ifeq ($(GOOS),windows)
  BIN 		:= watergun.exe
  DIST 		:= watergun_$(VERSION)_$(GOOS)_$(GOARCH).zip
else
  BIN		:= watergun
  DIST 		:= watergun_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz
endif

.PHONY: all build dist test clean

all: build

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-o "./$(TMP_DIR)/$(BIN)" \
		-ldflags "\
		  -s -w \
		  -X 'github.com/snaztoz/watergun/version.Version=$(VERSION)' \
		  -X 'github.com/snaztoz/watergun/version.Commit=$(GIT_COMMIT)' \
		  -X 'github.com/snaztoz/watergun/version.Date=$(BUILD_DATE)'" \
		"./$(MAIN_FILE)"

dist: build
	mkdir -p ./$(DIST_DIR)
	cp ./$(TMP_DIR)/$(BIN) ./$(DIST_DIR)/
	cp -r LICENSES ./$(DIST_DIR)/
	cp LICENSE NOTICE ./$(DIST_DIR)/

ifeq ($(GOOS),windows)
	powershell -NoProfile -Command \
	  "Compress-Archive -Force \
	    '$(DIST_DIR)\$(BIN)', \
	    '$(DIST_DIR)\LICENSES', \
	    '$(DIST_DIR)\LICENSE', \
	    '$(DIST_DIR)\NOTICE' \
	    '$(DIST_DIR)\$(DIST)'"
else
	tar -czf $(DIST_DIR)/$(DIST) -C $(DIST_DIR) \
	  $(BIN) LICENSES LICENSE NOTICE
endif

test:
	go test "./..."

clean:
	rm -rf ./$(DIST_DIR)/*

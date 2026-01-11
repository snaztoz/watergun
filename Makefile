VERSION 	?= dev
GIT_COMMIT 	:= $(shell git rev-parse --short HEAD)
BUILD_DATE 	:= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

DIST_DIR 	:= dist
TMP_DIR		:= tmp
MAIN_FILE 	:= watergun.go
BIN			:= watergun
DIST 		:= watergun-$(VERSION).tar.gz

.PHONY: all build dist test clean

all: build

build:
	go build \
		-o "./$(TMP_DIR)/$(BIN)" \
		-ldflags "\
			-X 'github.com/snaztoz/watergun/version.Version=$(VERSION)' \
			-X 'github.com/snaztoz/watergun/version.Commit=$(GIT_COMMIT)' \
			-X 'github.com/snaztoz/watergun/version.Date=$(BUILD_DATE)'" \
		"./$(MAIN_FILE)"

dist: build
	cp -r ./$(TMP_DIR)/$(BIN) LICENSES LICENSE NOTICE ./$(DIST_DIR)
	tar -czf ./$(DIST_DIR)/$(DIST) -C ./$(DIST_DIR)\
		$(BIN) \
		LICENSES \
		LICENSE \
		NOTICE

test:
	go test "./..."

clean:
	find ./$(DIST_DIR) -mindepth 1 ! -name '.gitignore' -delete

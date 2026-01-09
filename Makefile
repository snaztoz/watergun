OUT_DIR 	= ./tmp
MAIN_FILE 	= watergun.go
BUILD_OUT 	= watergun
DIST 		= watergun.tar.gz

.PHONY: all build dist test clean

all: build

build: clean
	go build -o "$(OUT_DIR)/$(BUILD_OUT)" "./$(MAIN_FILE)"

dist: build
	cp -r LICENSES LICENSE NOTICE $(OUT_DIR)
	tar -czf $(OUT_DIR)/$(DIST) -C $(OUT_DIR)\
		$(BUILD_OUT) \
		LICENSES \
		LICENSE \
		NOTICE

test:
	go test "./..."

clean:
	rm -rf \
		$(OUT_DIR)/LICENSES \
		$(OUT_DIR)/LICENSE \
		$(OUT_DIR)/NOTICE \
		$(OUT_DIR)/$(BUILD_OUT) \
		$(OUT_DIR)/$(DIST)

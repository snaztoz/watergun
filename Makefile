MAIN_FILE = "./watergun.go"
BUILD_OUT = "./tmp/watergun"

.PHONY: all build build-all test clean

all: build

build: clean
	go build -o $(BUILD_OUT) $(MAIN_FILE)

test:
	go test "./..."

clean:
	rm -f $(BUILD_OUT)

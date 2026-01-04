MAIN_FILE = "watergun.go"
OUT_FILE = "./tmp/watergun"

.PHONY: build test clean

build:
	go build -o $(OUT_FILE) $(MAIN_FILE)

test:
	go test "./..."

clean:
	rm $(OUT_FILE)

MAIN_FILE = "watergun.go"
OUT_FILE = "./tmp/watergun"

.PHONY: build test clean

build: clean
	go build -o $(OUT_FILE) $(MAIN_FILE)

test:
	go test "./..."

clean:
	rm -f $(OUT_FILE)

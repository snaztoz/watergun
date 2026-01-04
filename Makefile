MAIN_FILE = "watergun.go"
OUT_FILE = watergun

.PHONY: build clean

build:
	go build -o $(OUT_FILE) $(MAIN_FILE)

test:
	go test "./..."

clean:
	rm $(OUT_FILE)

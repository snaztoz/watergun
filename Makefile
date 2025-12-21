MAIN_FILE = "cmd/main.go"
OUT_FILE = watergun

.PHONY: build clean

build:
	go build -o $(OUT_FILE) $(MAIN_FILE)

clean:
	rm $(OUT_FILE)

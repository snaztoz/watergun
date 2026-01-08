SERVER_FILE = "./cmd/watergun/watergun.go"
SENDER_FILE = "./cmd/sender/sender.go"
LISTENER_FILE = "./cmd/listener/listener.go"

SERVER_OUT = "./tmp/watergun"
SENDER_OUT = "./tmp/watergun-send"
LISTENER_OUT = "./tmp/watergun-listen"

.PHONY: all build build-all test clean

all: build

build: clean
	go build -o $(SERVER_OUT) $(SERVER_FILE)

build-all: clean
	go build -o $(SERVER_OUT) $(SERVER_FILE)
	go build -o $(SENDER_OUT) $(SENDER_FILE)
	go build -o $(LISTENER_OUT) $(LISTENER_FILE)

test:
	go test "./..."

clean:
	rm -f $(SERVER_OUT) $(SENDER_OUT) $(LISTENER_OUT)

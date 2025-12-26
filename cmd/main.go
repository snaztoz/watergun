package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/snaztoz/watergun/server"
)

const port = "8080"

func main() {
	s := server.New(port)

	go s.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.Stop()
}

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/snaztoz/watergun/server"
)

const port = "8080"

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server process",
		Run: func(cmd *cobra.Command, args []string) {
			s := server.New(port)

			go s.Run()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			<-quit
			s.Stop()
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

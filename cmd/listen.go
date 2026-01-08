package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Listen for incoming messages inside a room",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listening for incoming messages")
		},
	}
)

func init() {
	rootCmd.AddCommand(listenCmd)
}

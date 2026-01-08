package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Quickly send message(s) to a room",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Sending message")
		},
	}
)

func init() {
	rootCmd.AddCommand(sendCmd)
}

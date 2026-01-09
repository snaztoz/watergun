package cmd

import (
	"fmt"

	"github.com/snaztoz/watergun/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Watergun %s (%s), compiled at %s\n", version.Version, version.Commit, version.Date)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

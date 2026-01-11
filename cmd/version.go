package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/snaztoz/watergun/version"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Watergun %s-%s-%s (%s, compiled at %s)\n", runtime.GOOS, runtime.GOARCH, version.Version, version.Commit, version.Date)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

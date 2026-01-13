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
		Short: "Display version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("watergun %s-%s %s/%s BuildDate=%s\n", version.Version, version.Commit, runtime.GOOS, runtime.GOARCH, version.Date)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

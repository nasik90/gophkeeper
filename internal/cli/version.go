package cli

import (
	"fmt"

	"github.com/nasik90/gophkeeper/internal/client/app"
	"github.com/spf13/cobra"
)

func VersionCommand() *cobra.Command {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version and build date",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nBuild Date: %s\n", app.Version, app.BuildDate)
		},
	}
	return versionCmd
}

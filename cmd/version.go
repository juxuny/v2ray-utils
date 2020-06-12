package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number of v2ray-cli",
	Long:  `This command can be used get the version number of v2ray-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v2ray-cli v0.0.1-alpha")
	},
}

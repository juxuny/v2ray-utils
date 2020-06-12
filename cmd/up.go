package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "update v2ray subscription",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("updating")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

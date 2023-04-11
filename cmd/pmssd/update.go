package pmssd

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update database",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

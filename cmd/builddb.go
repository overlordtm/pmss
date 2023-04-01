/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/overlordtm/pmss/pkg/sigimport"
	"github.com/spf13/cobra"
)

// builddbCmd represents the builddb command
var builddbCmd = &cobra.Command{
	Use:   "builddb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := sigimport.Import("test2.db", "data/full2.csv")
		if err != nil {
			return fmt.Errorf("error while loading signatures: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(builddbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	builddbCmd.PersistentFlags().String("dbPath", "pmss.db", "Path to the database file to be created")
	builddbCmd.PersistentFlags().String("csvPath", "data/full.csv", "Path to the csv file containing the signatures")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// builddbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

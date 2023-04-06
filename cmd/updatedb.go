/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/debscraper"
	"github.com/spf13/cobra"
)

// updatedbCmd represents the updatedb command
var updatedbCmd = &cobra.Command{
	Use:   "updatedb",
	Short: "A brief description of your command",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("updatedb called")

		scraper := debscraper.New()

		store, err := datastore.New(datastore.WithDbUrl(rootFlags.DBPath))
		if err != nil {
			return fmt.Errorf("error while creating datastore: %v", err)
		}

		packages, err := scraper.ListPackages(context.Background())

		if err != nil {
			return err
		}

		packageRows := make([]datastore.Package, len(packages))

		for i, pkg := range packages {
			packageRows[i].Name = pkg.Name
			packageRows[i].Version = pkg.Version
			packageRows[i].Architecture = pkg.Architecture
			packageRows[i].Filename = pkg.Filename
			packageRows[i].Size = pkg.Size
			packageRows[i].MD5 = pkg.MD5Sum
			packageRows[i].SHA256 = pkg.SHA256
		}

		store.Packages().InsertBatch(packageRows)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updatedbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updatedbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updatedbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"strings"

	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/datastore/sqlitestore"
	"github.com/overlordtm/pmss/pkg/debscraper"
	"github.com/spf13/cobra"
)

var (
	dbPath   string
	filePath string
)

// builddbCmd represents the builddb command
var builddbCmd = &cobra.Command{
	Use:   "builddb",
	Short: "Imports data from a file into database",
	Long:  `Imports data from a file into database`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		dialector, err := datastore.ParseDBUrl(dbPath)
		if err != nil {
			return err
		}

		var ds *datastore.Store

		if strings.HasSuffix(dbPath, ".sqlite3") {
			ds, err = datastore.New(datastore.WithDb(dialector))
			if err != nil {
				return fmt.Errorf("error while creating hash store: %v", err)
			}
		} else {
			return fmt.Errorf("invalid database file type")
		}

		fileName := filepath.Base(filePath)

		if strings.HasSuffix(fileName, ".csv") {
			// we are dealing with csv file

			if strings.HasPrefix(fileName, "deb") {
				//
				file, err := os.Open(filePath)
				if err != nil {
					return fmt.Errorf("error while opening file: %v", err)
				}
				decoder, err := debscraper.NewCsvDecoder(file)
				if err != nil {
					return fmt.Errorf("error while creating csv decoder: %v", err)
				}

				batch := make([]datastore.WhitelistItem, 0, sqlitestore.WhitelistBatchSize)
				for {
					item := debscraper.HashItem{}
					err := decoder.Decode(&item)
					if err == io.EOF {
						break
					}
					if err != nil {
						return fmt.Errorf("error while decoding csv: %v", err)
					}

					batch = append(batch, datastore.WhitelistItem{
						MD5:    item.MD5,
						SHA1:   item.SHA1,
						SHA256: item.SHA256,
						Path:   item.Filename,
						Meta: datastore.WhitelistMeta{
							Package: item.Package,
							Version: item.Version,
							Size:    item.Size,
							Owner:   item.Owner,
							Group:   item.Group,
							Mode:    uint32(item.Mode),
						},
					})

					if len(batch) == sqlitestore.WhitelistBatchSize {
						err := ds.Whitelist().InsertBatch(batch)
						if err != nil {
							return fmt.Errorf("error while adding to whitelist: %v", err)
						}
						batch = batch[:0]
					}
				}

				if len(batch) > 0 {
					for _, item := range batch {
						err := ds.Whitelist().Insert(item)
						if err != nil {
							return fmt.Errorf("error while adding to whitelist: %v", err)
						}
					}
				}

			} else {
				return fmt.Errorf("unknown file type %s", fileName)
			}

		} else {
			return fmt.Errorf("unknown file extension")
		}

		// err := sigimport.Import(dbPath, filePath)
		// if err != nil {
		// 	return fmt.Errorf("error while loading signatures: %v", err)
		// }
		// return nil
		return nil
	},
}

func init() {
	rootCmd.AddCommand(builddbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	builddbCmd.PersistentFlags().StringVar(&dbPath, "dbPath", "pmss.sqlite3", "Path to the database file to be created")
	builddbCmd.PersistentFlags().StringVar(&filePath, "csvPath", "data/full.csv", "Path to the csv file containing the signatures")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// builddbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

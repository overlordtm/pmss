/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/overlordtm/pmss/pkg/debscraper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	workers        int
	url            string
	outputFilePath string
	distro         string
	arch           string
	component      string
)

// debscrapeCmd represents the debscrape command
var debscrapeCmd = &cobra.Command{
	Use:   "debscrape",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if outputFilePath == "" {
			outputFilePath = fmt.Sprintf("debscrape-%s-%s-%s-%s.csv", distro, arch, component, time.Now().Format("2006-01-02-15-04-05"))
		}
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		logrus.SetOutput(os.Stderr)

		scraper := debscraper.New(debscraper.WithMirrorUrl(url), debscraper.WithDistro(distro), debscraper.WithArch(arch), debscraper.WithComponent(component))

		results := make(chan debscraper.HashItem, 1024)

		go func() {
			err = scraper.Scrape(workers, results)
		}()

		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		csvWriter := csv.NewWriter(outFile)
		defer csvWriter.Flush()
		encoder := csvutil.NewEncoder(csvWriter)

		// cvsWriter.Write([]string{"filename", "package", "version", "architecture", "md5", "sha1", "sha256", "size", "mode", "owner", "group"})

		for result := range results {
			err = encoder.Encode(result)
			// err = cvsWriter.Write([]string{result.Filename, result.Package, result.Version, result.Architecture, result.MD5, result.SHA1, result.SHA256, fmt.Sprintf("%d", result.Size), fmt.Sprintf("%o", result.Mode), result.Owner, result.Group})
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(debscrapeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// debscrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// debscrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	debscrapeCmd.Flags().IntVarP(&workers, "workers", "w", 10, "Number of workers")
	debscrapeCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Output file")
	debscrapeCmd.Flags().StringVarP(&url, "url", "u", "http://ftp.debian.org/debian", "Mirror URL")
	debscrapeCmd.Flags().StringVarP(&arch, "arch", "a", "amd64", "Architecture")
	debscrapeCmd.Flags().StringVarP(&distro, "distro", "d", "buster", "Distribution")
	debscrapeCmd.Flags().StringVarP(&component, "component", "c", "main", "Component")
}

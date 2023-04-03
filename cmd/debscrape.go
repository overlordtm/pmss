/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	logfile        string
	progress       bool = true
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

		opts := make([]debscraper.Option, 0)

		if logfile != "" {
			logger := logrus.New()
			logger.SetFormatter(&logrus.TextFormatter{
				DisableColors: true,
				FullTimestamp: true,
			})
			logger.SetLevel(logrus.DebugLevel)

			logFileFd, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("error while opening log file: %w", err)
			}
			defer logFileFd.Close()

			logger.SetOutput(logFileFd)

			// flush logfile every 5 seconds
			// ticker := time.NewTicker(5 * time.Second)
			// defer ticker.Stop()
			// go func() {
			// 	for {
			// 		select {
			// 		case <-ticker.C:
			// 			logFileFd.Sync()
			// 		}
			// 	}
			// }()

			opts = append(opts, debscraper.WithLogger(logger))

		} else {
			logrus.SetOutput(os.Stdout)
		}

		if url != "" {
			opts = append(opts, debscraper.WithMirrorUrl(url))
		}

		if distro != "" {
			opts = append(opts, debscraper.WithDistro(distro))
		}

		if arch != "" {
			opts = append(opts, debscraper.WithArch(arch))
		}

		if component != "" {
			opts = append(opts, debscraper.WithComponent(component))
		}

		scraper := debscraper.New(opts...)

		results := make(chan debscraper.HashItem, 1024)

		// // print some runtime stats
		// go func() {
		// 	ticker := time.NewTicker(5 * time.Second)
		// 	for {
		// 		select {
		// 		case <-ticker.C:
		// 			memStats := runtime.MemStats{}
		// 			runtime.ReadMemStats(&memStats)
		// 			logrus.Infof("gorutines: %d, mem: %s, mallocs: %d", runtime.NumGoroutine(), humanize.Bytes(memStats.Alloc), memStats.Mallocs)
		// 		}
		// 	}
		// }()

		ctx, cancel := context.WithCancel(context.Background())

		// listen for SIGINT and SIGTERM and call cancel on the context
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
			<-sigint
			cancel()
			time.Sleep(200 * time.Millisecond)
			os.Exit(1)
		}()

		go func() {
			var pbDelegate debscraper.ProgressDelegate
			if progress {
				pbDelegate = &debscraper.CliProgressBar{}
			} else {
				pbDelegate = &debscraper.NoopProgressBar{}
			}

			err = scraper.Scrape(ctx, workers, results, pbDelegate)
		}()

		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer outFile.Close()

		csvWriter := csv.NewWriter(outFile)
		defer csvWriter.Flush()
		// encoder := csvutil.NewEncoder(csvWriter)
		encoder := debscraper.NewCsvEncoder(outFile)

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
	debscrapeCmd.Flags().StringVarP(&url, "url", "u", "", "Mirror URL (if empty, default list is used in round robim node)")
	debscrapeCmd.Flags().StringVarP(&arch, "arch", "a", "amd64", "Architecture")
	debscrapeCmd.Flags().StringVarP(&distro, "distro", "d", "buster", "Distribution")
	debscrapeCmd.Flags().StringVarP(&component, "component", "c", "main", "Component")
	debscrapeCmd.Flags().BoolVarP(&progress, "progress", "p", false, "Show progress bar")
	debscrapeCmd.Flags().StringVar(&logfile, "log", "", "Log file")
}

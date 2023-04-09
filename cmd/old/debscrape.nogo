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
	"runtime"
	"syscall"
	"time"

	"github.com/overlordtm/pmss/pkg/debscraper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type DebscrapeFlags struct {
	Workers        int
	Url            string
	OutputFilePath string
	Distro         string
	Arch           string
	Component      string
	Logfile        string
	Progress       bool
}

var debscrapeFlags DebscrapeFlags = DebscrapeFlags{
	Workers:        runtime.NumCPU() * 2,
	Url:            "",
	OutputFilePath: "",
	Distro:         "buster",
	Arch:           "amd64",
	Component:      "main",
	Logfile:        "",
	Progress:       false,
}

func init() {
	debscrapeFlags.OutputFilePath = fmt.Sprintf("debscrape-%s-%s-%s-%s.csv", debscrapeFlags.Distro, debscrapeFlags.Arch, debscrapeFlags.Component, time.Now().Format("2006-01-02-15-04-05"))
}

// debscrapeCmd represents the debscrape command
var debscrapeCmd = &cobra.Command{
	Use:   "debscrape",
	Short: "Scrapes packages in debian repository for hashes",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		opts := make([]debscraper.Option, 0)

		if debscrapeFlags.Logfile != "" {
			logger := logrus.New()
			logger.SetFormatter(&logrus.TextFormatter{
				DisableColors: true,
				FullTimestamp: true,
			})

			logger.SetLevel(logrus.GetLevel())

			logFileFd, err := os.OpenFile(debscrapeFlags.Logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("error while opening log file: %w", err)
			}
			defer logFileFd.Close()

			logger.SetOutput(logFileFd)

			opts = append(opts, debscraper.WithLogger(logger))

		} else {
			logrus.SetOutput(os.Stdout)
		}

		if debscrapeFlags.Url != "" {
			opts = append(opts, debscraper.WithMirrorUrl(debscrapeFlags.Url))
		}

		if debscrapeFlags.Distro != "" {
			opts = append(opts, debscraper.WithDistro(debscrapeFlags.Distro))
		}

		if debscrapeFlags.Arch != "" {
			opts = append(opts, debscraper.WithArch(debscrapeFlags.Arch))
		}

		if debscrapeFlags.Component != "" {
			opts = append(opts, debscraper.WithComponent(debscrapeFlags.Component))
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
			if debscrapeFlags.Progress {
				pbDelegate = &debscraper.CliProgressBar{}
			} else {
				pbDelegate = &debscraper.NoopProgressBar{}
			}

			err = scraper.Scrape(ctx, debscrapeFlags.Workers, results, pbDelegate)
		}()

		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(debscrapeFlags.OutputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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
	debscrapeCmd.Flags().IntVar(&debscrapeFlags.Workers, "workers", debscrapeFlags.Workers, "Number of workers")
	debscrapeCmd.Flags().StringVar(&debscrapeFlags.OutputFilePath, "output", debscrapeFlags.OutputFilePath, "Output file")
	debscrapeCmd.Flags().StringVar(&debscrapeFlags.Url, "url", "", "Mirror URL (if empty, default list is used in round robim node)")
	debscrapeCmd.Flags().StringVar(&debscrapeFlags.Arch, "arch", "amd64", "Architecture")
	debscrapeCmd.Flags().StringVar(&debscrapeFlags.Distro, "distro", "buster", "Distribution")
	debscrapeCmd.Flags().StringVar(&debscrapeFlags.Component, "component", "main", "Component")
	debscrapeCmd.Flags().BoolVar(&debscrapeFlags.Progress, "progress", false, "Show progress bar")
	debscrapeCmd.Flags().StringVar(&debscrapeFlags.Logfile, "log", "", "Log file")
}

package pmss

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/overlordtm/pmss/internal/apiclient"
	"github.com/overlordtm/pmss/pkg/client"
	"github.com/overlordtm/pmss/pkg/scanner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagApiUrl       = "api-url"
	flagPathsExclude = "exclude"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans given list of files or directoryes",
	Long:  `Scans given list of files or directoryes, compute hashes and extract metadata and report it to the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		logrus.SetOutput(os.Stderr)

		apiUrl := viper.GetString("api.url")
		excludePaths := viper.GetStringSlice("paths.exclude")

		ctx, _ := context.WithCancel(context.Background())

		client, err := client.New(apiUrl)
		if err != nil {
			return err
		}

		// default to scanning root
		if len(args) == 0 {
			args = append(args, "/")
		}

		// convert paths in args to absolute paths
		for i, arg := range args {
			absPath, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %v", arg, err)
			}
			args[i] = absPath
		}

		scn := scanner.New(scanner.WithExcludePaths(excludePaths))

		ch := make(chan apiclient.File, 1024)

		go func() {
			err := scn.Scan(ch, args...)
			if err != nil {
				logrus.WithError(err).Error("error while scanning")
			}
		}()

		batchSize := 1000

		files := make([]apiclient.File, 0, batchSize)

		var reportRunId *uuid.UUID = nil

		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("failed to get hostname: %v", err)
		}

		machineId, err := machineid.ProtectedID("pmss")
		if err != nil {
			return fmt.Errorf("failed to get machine id: %v", err)
		}

		for {
			select {
			case f, ok := <-ch:

				if ok {
					logrus.WithField("file", f.Path).Debug("scanning file")
					files = append(files, f)
				}

				if len(files) == batchSize || !ok {
					logrus.WithField("files", len(files)).Debug("Submiting report batch")
					response, err := client.SubmitFiles(ctx, apiclient.NewReportRequest{
						Files:       files,
						Hostname:    hostname,
						MachineId:   machineId,
						ReportRunId: reportRunId,
					})
					if err != nil {
						logrus.WithError(err).Error("Failed to submit report batch")
					}

					if response.StatusCode() == http.StatusCreated {
						reportRunId = &response.JSON201.Id
						for _, file := range response.JSON201.Files {
							fmt.Printf("%s\t%s\n", file.Path, file.Status)
						}
					} else {
						logrus.WithField("statusCode", response.StatusCode()).Error("Failed to submit report batch, unexpected status code")
						return fmt.Errorf("failed to send files: %#+v", response.JSONDefault)
					}
					files = files[:0]
				}
				if !ok {
					return nil
				}
			}
		}

		return nil
	},
}

func init() {

	viper.SetDefault("paths.exclude", []string{"/dev/", "/sys/", "/proc/"})

	scanCmd.Flags().String(flagApiUrl, "http://localhost:8080/api/v1", "API URL to send the report to")
	scanCmd.Flags().StringSlice(flagPathsExclude, []string{"/proc/", "/dev/", "/sys/"}, "Paths to exclude")

	viper.GetViper().BindPFlag("api.url", scanCmd.Flags().Lookup(flagApiUrl))
	viper.GetViper().BindPFlag("paths.exclude", scanCmd.Flags().Lookup(flagPathsExclude))

	rootCmd.AddCommand(scanCmd)
}

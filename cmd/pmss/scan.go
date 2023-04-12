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
	flagApiUrl = "api-url"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans given list of files or directoryes",
	Long:  `Scans given list of files or directoryes, compute hashes and extract metadata and report it to the server.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		logrus.SetOutput(os.Stderr)

		apiUrl := viper.GetString("api.url")

		ctx, _ := context.WithCancel(context.Background())

		client, err := client.New(apiUrl)
		if err != nil {
			return err
		}

		// convert paths in args to absolute paths
		for i, arg := range args {
			absPath, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %v", arg, err)
			}
			args[i] = absPath
		}

		scn := scanner.New()

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

		for f := range ch {
			logrus.WithField("file", f.Path).WithField("hash", f.Md5).Info("scanning file")
			files = append(files, f)

			if len(files) == batchSize {
				logrus.WithField("files", len(files)).Info("sending batch")
				response, err := client.SubmitFiles(ctx, apiclient.NewReportRequest{
					Files:       files,
					Hostname:    hostname,
					MachineId:   machineId,
					ReportRunId: reportRunId,
				})
				if err != nil {
					logrus.WithError(err).Error("failed to send files")
				}

				if response.StatusCode() == http.StatusCreated {
					reportRunId = &response.JSON201.Id
					for _, file := range response.JSON201.Files {
						fmt.Printf("%s\t%s\n", file.Path, file.Status)
					}
				} else {
					logrus.WithField("statusCode", response.StatusCode()).Error("failed to send files, unexpected status code")
					return fmt.Errorf("failed to send files: %#+v", response.JSONDefault)
				}
				files = files[:0]
			}
		}

		return nil
	},
}

func init() {
	viper.SetDefault("api.url", "http://localhost:8080/api/v1")

	scanCmd.Flags().String(flagApiUrl, viper.GetString("api-url"), "API URL to send the report to")

	viper.GetViper().BindPFlag("api.url", scanCmd.Flags().Lookup("api-url"))

	rootCmd.AddCommand(scanCmd)
}

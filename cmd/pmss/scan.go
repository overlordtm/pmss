package pmss

import (
	"context"
	"fmt"
	"net/http"

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

		apiUrl := viper.GetString("api.url")

		ctx, _ := context.WithCancel(context.Background())

		client, err := client.New(apiUrl)
		if err != nil {
			return err
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

		for f := range ch {
			files = append(files, f)

			if len(files) == batchSize {
				logrus.WithField("files", len(files)).Info("sending batch")
				response, err := client.SubmitFiles(ctx, apiclient.NewReportRequest{
					Files:       files,
					Hostname:    "test",
					MachineId:   "test",
					ReportRunId: reportRunId,
				})
				if err != nil {
					logrus.WithError(err).Error("failed to send files")
				}

				if response.StatusCode() == http.StatusCreated {
					// copy pointer value
					var tmpUuid = response.JSON201.Id
					reportRunId = &tmpUuid

				} else {
					return fmt.Errorf("failed to send files: %#+v", response.JSONDefault)
				}
				files = files[:0]
			}

			// logrus.WithField("file", f).Info("file scanned")
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

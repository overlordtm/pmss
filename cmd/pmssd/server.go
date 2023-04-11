package pmssd

import (
	"context"
	"fmt"

	"github.com/overlordtm/pmss/pkg/httpserver"
	"github.com/overlordtm/pmss/pkg/pmss"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagHttpListenAddr = "http-listen-addr"
	flagDbUrl          = "db-url"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server instance",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		httpListenAddr := viper.GetString("http.listen.address")
		dbUrl := viper.GetString("db.url")

		pmss, err := pmss.New(pmss.WithDbUrl(dbUrl))
		if err != nil {
			return fmt.Errorf("failed to initialize PMSS: %v", err)
		}

		logrus.WithField("httpListenAddr", httpListenAddr).Info("Starting HTTP server")
		srv := httpserver.New(context.Background(), pmss, httpserver.WithListenAddr(httpListenAddr))
		return srv.Start()
	},
}

func init() {
	viper.SetDefault("http.listen.address", ":8080")
	viper.SetDefault("db.url", "sqlite3://:memory:")

	viper.SetEnvPrefix("PMSS")
	viper.AutomaticEnv()

	serverCmd.Flags().String(flagHttpListenAddr, ":8080", "HTTP listen address")
	serverCmd.Flags().String(flagDbUrl, "sqlite3://:memory:", "Database URL")
	viper.BindPFlag("http.listen.address", serverCmd.Flags().Lookup(flagHttpListenAddr))
	viper.BindPFlag("db.url", serverCmd.Flags().Lookup(flagDbUrl))

	rootCmd.AddCommand(serverCmd)
}

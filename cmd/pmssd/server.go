package pmssd

import (
	"context"
	"fmt"
	"strings"

	"github.com/overlordtm/pmss/internal/utils"
	"github.com/overlordtm/pmss/internal/viperutils"
	"github.com/overlordtm/pmss/pkg/httpserver"
	"github.com/overlordtm/pmss/pkg/pmss"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagHttpListenAddr = "http.listen.addr"
	flagDbUrl          = "db.url"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server instance",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		httpListenAddr := viper.GetString(flagHttpListenAddr)
		dbUrl := viper.GetString(flagDbUrl)

		logrus.WithField("dbUrl", dbUrl).Info("Initializing PMSS")

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
	serverCmd.Flags().String(utils.ReplaceDotWithDash(flagHttpListenAddr), ":8080", "HTTP listen address")
	serverCmd.Flags().String(utils.ReplaceDotWithDash(flagDbUrl), "mysql://pmss:pmss@tcp(localhost:3306)/pmss?charset=utf8&parseTime=True&loc=Local", "Database URL")
	viper.BindPFlag(flagHttpListenAddr, serverCmd.Flags().Lookup(utils.ReplaceDotWithDash(flagHttpListenAddr)))
	viper.BindPFlag(flagDbUrl, serverCmd.Flags().Lookup(utils.ReplaceDotWithDash(flagDbUrl)))
	viperutils.BindPFlags(serverCmd.Flags(), strings.NewReplacer("-", "."))
	rootCmd.AddCommand(serverCmd)
}

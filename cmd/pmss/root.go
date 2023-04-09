package pmss

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "pmss",
	Short: "Poor man security scanner",
	Long:  ``,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if viper.GetBool("verbose") {
			logrus.SetLevel(logrus.InfoLevel)
		}

		if viper.GetBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
			// run http pprof server
			go func() {
				go func() {
					addr := "localhost:6060"
					logrus.Debugf("Starting pprof server on %s", addr)
					logrus.Debugln(http.ListenAndServe(addr, nil))
				}()
			}()
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	cobra.OnInitialize(initViper)
}

package pmssd

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// type cliFlags struct {
// 	Verbose bool
// 	Debug   bool
// 	DB      struct {
// 		Url string
// 	}
// 	Http struct {
// 		Listen struct {
// 			Address string
// 		}
// 	}
// }

// var config = cliFlags{}

var rootCmd = &cobra.Command{
	Use:   "pmssd",
	Short: "Poor man security scanner",
	Long:  `Poor man sercurity scanner server side`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// if rootFlags.Verbose {
		// 	logrus.SetLevel(logrus.InfoLevel)
		// }

		// if rootFlags.Debug {
		// 	logrus.SetLevel(logrus.DebugLevel)

		// 	// run http pprof server
		// 	go func() {
		// 		go func() {
		// 			addr := "localhost:6060"
		// 			logrus.Debugf("Starting pprof server on %s", addr)
		// 			logrus.Debugln(http.ListenAndServe(addr, nil))
		// 		}()
		// 	}()
		// }
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {

	cobra.OnInitialize(func() {
		viper.SetConfigName("pmssd")              // name of config file (without extension)
		viper.AddConfigPath("/etc/pmss")          // path to look for the config file in
		viper.AddConfigPath("$HOME/.config/pmss") // call multiple times to add many search paths
		viper.AddConfigPath("$HOME/.pmss/")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				logrus.WithError(err).Fatal("Error reading config file")
			} else {
				logrus.Info("No config file found")
			}
		}
		logrus.Debug("Using config file: ", viper.ConfigFileUsed())

		// env bindings
		viper.AutomaticEnv()
		viper.SetEnvPrefix("PMSS")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	})

}

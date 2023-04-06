/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CliFlags struct {
	Verbose    bool
	Debug      bool
	DebugPprof bool
	LogLevel   string
	DBPath     string
}

var rootFlags CliFlags = CliFlags{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pmss",
	Short: "Poor man security scanner",
	Long:  ``,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if rootFlags.Verbose {
			logrus.SetLevel(logrus.InfoLevel)
		}

		if rootFlags.Debug {
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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pmss.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.Verbose, "verbose", "v", rootFlags.Verbose, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.Debug, "debug", "d", rootFlags.Debug, "debug output")
	rootCmd.PersistentFlags().StringVar(&rootFlags.DBPath, "db-path", "pmss.sqlite3", "path to database file")
	rootCmd.PersistentFlags().StringVar(&rootFlags.LogLevel, "log-level", rootFlags.LogLevel, "log level (debug, info, warn, error, fatal, panic)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

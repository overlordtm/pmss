/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/overlordtm/pmss/pkg/httpserver"
	"github.com/overlordtm/pmss/pkg/pmss"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ServerFlags struct {
	httpListenAddr string
}

var serverFlags = ServerFlags{}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server instance",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		pmss, err := pmss.New(rootFlags.DBUrl)
		if err != nil {
			return fmt.Errorf("failed to initialize PMSS: %v", err)
		}

		logrus.WithField("httpListenAddr", serverFlags.httpListenAddr).Info("Starting HTTP server")
		srv := httpserver.New(context.Background(), pmss, httpserver.WithListenAddr(serverFlags.httpListenAddr))
		return srv.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringVar(&serverFlags.httpListenAddr, "http-listen-addr", ":8080", "HTTP listen address")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/overlordtm/pmss/pkg/pmss"
	"github.com/overlordtm/pmss/pkg/server/httpserver"
	"github.com/spf13/cobra"
)

var (
	httpListenAddr string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("server called")

		pmss, err := pmss.New("test.db")
		if err != nil {
			return fmt.Errorf("failed to initialize PMSS: %v", err)
		}

		srv := httpserver.New(context.Background(), pmss)
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
	serverCmd.Flags().StringVar(&httpListenAddr, "http-listen-addr", ":8080", "HTTP listen address")
}

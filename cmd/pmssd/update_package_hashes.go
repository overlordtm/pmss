package pmssd

import (
	"context"
	"fmt"
	"runtime"

	"github.com/overlordtm/pmss/pkg/pmss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updatePackageHashesCmd = &cobra.Command{
	Use:   "packagehash",
	Short: "Scrape package's file hashes",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("update packages called")

		dbUrl := viper.GetString("db.url")

		pmss, err := pmss.New(pmss.WithDbUrl(dbUrl))
		if err != nil {
			return fmt.Errorf("failed to initialize PMSS: %v", err)
		}

		return pmss.UpdatePackageHashes(context.Background(), runtime.NumCPU()*2)
	},
}

func init() {
	updateCmd.AddCommand(updatePackageHashesCmd)
}

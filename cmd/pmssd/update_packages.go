package pmssd

import (
	"context"
	"fmt"

	"github.com/overlordtm/pmss/pkg/pmss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updatePackagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "Update packages database",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("update packages called")

		dbUrl := viper.GetString("db.url")

		pmss, err := pmss.New(pmss.WithDbUrl(dbUrl))
		if err != nil {
			return fmt.Errorf("failed to initialize PMSS: %v", err)
		}

		return pmss.UpdatePackages(context.Background())
	},
}

func init() {
	updateCmd.AddCommand(updatePackagesCmd)
}

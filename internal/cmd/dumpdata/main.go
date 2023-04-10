package main

import (
	"fmt"
	"os"

	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "dumpdata",
		RunE: func(cmd *cobra.Command, args []string) error {

			db := datastore.MustOpen("mysql://pmss:pmss@tcp(mariadb:3306)/test_pmss?parseTime=true")

			all := make([]datastore.KnownFile, 0)

			if err := datastore.KnownFiles().All(all)(db); err != nil {
				return err
			}

			fmt.Println(all)
			return nil
		},
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

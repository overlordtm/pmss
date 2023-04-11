package main

import (
	"os"

	"github.com/overlordtm/pmss/cmd/pmss"
)

func main() {
	if err := pmss.Execute(); err != nil {
		os.Exit(1)
	}
}

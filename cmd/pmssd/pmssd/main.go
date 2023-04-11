package main

import (
	"os"

	"github.com/overlordtm/pmss/cmd/pmssd"
)

func main() {
	if err := pmssd.Execute(); err != nil {
		os.Exit(1)
	}
}

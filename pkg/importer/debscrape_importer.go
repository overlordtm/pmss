package importer

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jszwec/csvutil"
	"github.com/overlordtm/pmss/pkg/debscraper"
)

func DebscrapeImporter(r io.Reader) error {
	decoder, err := csvutil.NewDecoder(csv.NewReader(r))
	if err != nil {
		return fmt.Errorf("error while creating decoder: %w", err)
	}

	row := &debscraper.HashItem{}

	for {
		if err := decoder.Decode(row); err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("error while decoding row: %w", err)
		}
		fmt.Printf("%+v\n", row)
	}
	return nil
}

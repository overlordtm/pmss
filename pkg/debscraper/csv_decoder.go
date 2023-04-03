package debscraper

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jszwec/csvutil"
)

type CsvDecoder struct {
	*csvutil.Decoder
}

func NewCsvDecoder(r io.Reader) (*CsvDecoder, error) {
	dec, err := csvutil.NewDecoder(csv.NewReader(r))
	if err != nil {
		return nil, fmt.Errorf("error while creating csv decoder: %v", err)
	}
	return &CsvDecoder{dec}, nil
}

func (c *CsvDecoder) Decode(v *HashItem) error {
	return c.Decoder.Decode(v)
}

package debscraper

import (
	"encoding/csv"
	"io"

	"github.com/jszwec/csvutil"
)

type CsvEncoder struct {
	*csvutil.Encoder
}

func NewCsvEncoder(w io.Writer) *CsvEncoder {
	return &CsvEncoder{csvutil.NewEncoder(csv.NewWriter(w))}
}

func (c *CsvEncoder) Encode(v HashItem) error {
	return c.Encoder.Encode(v)
}

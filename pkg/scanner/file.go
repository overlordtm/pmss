package scanner

import (
	"fmt"
	"os"

	"errors"

	"github.com/overlordtm/pmss/pkg/checker"
	"github.com/overlordtm/pmss/pkg/multihasher"
	"github.com/overlordtm/pmss/pkg/sigdb"
)

func scanFile(db sigdb.SigDb, path string, checkers []interface{}) (r Result, err error) {

	r.Path = path

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return r, fmt.Errorf("error while opening file: %v", err)
	}

	defer f.Close()

	h, err := multihasher.Hash(f)

	for _, c := range checkers {

		switch c1 := c.(type) {
		case checker.HashChecker:
			res, err1 := c1.CheckHash(h)
			if err != nil {
				err = errors.Join(err, err1)
				return r, fmt.Errorf("error while checking hash: %v", err)
			}

			if res != nil {
				r.CheckResults = append(r.CheckResults, *res)
			}

		default:
			return r, fmt.Errorf("unknown checker type")
		}
	}

	return r, nil
}

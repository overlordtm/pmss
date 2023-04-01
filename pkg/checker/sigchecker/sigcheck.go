package sigchecker

import (
	"errors"

	"github.com/overlordtm/pmss/pkg/checker"
	"github.com/overlordtm/pmss/pkg/multihasher"
	"github.com/overlordtm/pmss/pkg/sigdb"
)

type SigChecker struct {
	db sigdb.SigDb
}

func New(db sigdb.SigDb) *SigChecker {
	return &SigChecker{
		db: db,
	}
}

func (c *SigChecker) CheckHash(h multihasher.Hashes) (r *checker.Result, err2 error) {

	if h.MD5 != "" {
		item, err := c.db.FindByMD5(h.MD5)

		if err != nil {
			err2 = errors.Join(err2, err)
		}

		if item != nil {
			return &checker.Result{
				Signature: item.Signature,
			}, err2
		}
	}

	if h.SHA1 != "" {
		item, err := c.db.FindBySHA1(h.SHA1)

		if err != nil {
			err2 = errors.Join(err2, err)
		}

		if item != nil {
			return &checker.Result{
				Signature: item.Signature,
			}, err2
		}
	}

	if h.SHA256 != "" {
		item, err := c.db.FindBySHA256(h.SHA256)

		if err != nil {
			err2 = errors.Join(err2, err)
		}

		if item != nil {
			return &checker.Result{
				Signature: item.Signature,
			}, err2
		}
	}

	return nil, err2
}

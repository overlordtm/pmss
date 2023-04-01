package pmss

import (
	"github.com/overlordtm/pmss/pkg/hashvariant"
	"github.com/overlordtm/pmss/pkg/sigdb"
)

type Pmss struct {
	db sigdb.SigDb
}

type Result struct {
	Hash        string
	HashVariant hashvariant.HashVariant
	KnownPaths  []string
}

func New(dbPath string) (*Pmss, error) {

	return &Pmss{}, nil
}

func (p *Pmss) FindByHash(hash string) (r *Result, err error) {

	r = &Result{}

	variant := hashvariant.DetectHashVariant(hash)

	r.HashVariant = variant

	return r, nil
}

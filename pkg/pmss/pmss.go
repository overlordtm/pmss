package pmss

import (
	"fmt"

	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/hashvariant"
)

type Pmss struct {
	ds *datastore.Store
}

type Result struct {
	Hash        string
	HashVariant hashvariant.HashVariant
	KnownPaths  []string
}

func New(dbPath string) (*Pmss, error) {

	dialector, err := datastore.ParseDBUrl(dbPath)
	if err != nil {
		return nil, err
	}

	ds, err := datastore.New(datastore.WithDb(dialector))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize datastore: %v", err)
	}

	return &Pmss{
		ds: ds,
	}, nil
}

func (p *Pmss) FindByHash(hash string) (r *Result, err error) {

	r = &Result{}

	variant := hashvariant.DetectHashVariant(hash)

	r.HashVariant = variant

	return r, nil
}

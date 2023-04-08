package pmss

import (
	"fmt"

	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/hashvariant"
)

type Pmss struct {
	Data *datastore.Store
}

type HashSearchResult struct {
	Hash        string                  `json:"hash"`
	HashVariant hashvariant.HashVariant `json:"hash_variant"`
	KnownFiles  []string                `json:"known_files"`
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
		Data: ds,
	}, nil
}

func (p *Pmss) FindByHash(hash string) (malformed bool, r *HashSearchResult, err error) {
	var files []string
	t := hashvariant.HashVariant("unknown")
	err = p.Data.KnownFiles().FindAllPathsByHash(hash, &files, &t)
	malformed = false
	if t == hashvariant.Unknown {
		malformed = true
		t = hashvariant.HashVariant("unknown")
	}
	r = &HashSearchResult{
		Hash:        hash,
		HashVariant: t,
		KnownFiles:  files,
	}
	return
}

func (p *Pmss) FindMachineByHostname(machineHostname, machineApiKey string, machine *datastore.Machine) (bool, error) {
	if err := p.Data.Machines().FindByHostname(machineHostname, machine); err != nil {
		return false, err
	}
	if machine.ApiKey != machineApiKey {
		return false, nil
	}
	return true, nil
}

func (p *Pmss) DoMachineReport(reports []datastore.ScannedFile, machine *datastore.Machine, dest *datastore.ReportRun) error {

	if err := p.Data.ReportRuns().CreateNew(dest); err != nil {
		return err
	}

	for _, d := range reports {
		d.ReportRun = *dest
		d.Machine = *machine
	}

	if err := p.Data.ScannedFiles().InsertBatch(reports); err != nil {
		return err
	}

	return nil
}

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

func (p *Pmss) FindByHash(hash string) (r *HashSearchResult, err error) {
	var files []string

	t := hashvariant.DetectHashVariant(hash)

	if t == hashvariant.Unknown {
		return nil, fmt.Errorf("unknown hash variant")
	}

	err = p.Data.KnownFiles().FindAllPathsByHash(hash, &files, &t)
	if err != nil {
		return nil, err
	}
	return &HashSearchResult{
		Hash:        hash,
		HashVariant: t,
		KnownFiles:  files,
	}, nil
}

func (p *Pmss) FindMachineByHostname(machineHostname string, machine *datastore.Machine) (bool, error) {
	if err := p.Data.Machines().FindByHostname(machineHostname, machine); err != nil {
		return false, err
	}
	// if machine.ApiKey != machineApiKey {
	// 	return false, nil
	// }
	return true, nil
}

func (p *Pmss) DoMachineReport(files []datastore.ScannedFile, machine *datastore.Machine, reportRun *datastore.ReportRun) error {

	if err := p.Data.ReportRuns().CreateNew(reportRun); err != nil {
		return err
	}

	for _, d := range files {
		d.ReportRun = *reportRun
		d.Machine = *machine
	}

	if err := p.Data.ScannedFiles().InsertBatch(files); err != nil {
		return err
	}

	return nil
}

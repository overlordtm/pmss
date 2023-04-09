package pmss

import (
	"fmt"

	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/hashvariant"
	"gorm.io/gorm"
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

func (p *Pmss) DoMachineReport(scanReport *ScanReport) (*datastore.ReportRun, error) {

	var reportRun *datastore.ReportRun

	if scanReport.ScanRunId != nil {
		if err := p.Data.ReportRuns().FindByID(*scanReport.ScanRunId, reportRun); err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find report run: %v", err)
		}
	}

	if reportRun == nil {
		reportRun = &datastore.ReportRun{}
		if err := p.Data.ReportRuns().CreateNew(reportRun); err != nil {
			return nil, fmt.Errorf("failed to create new report run: %v", err)
		}
	}

	var machine *datastore.Machine = &datastore.Machine{
		MachineId: scanReport.MachineId,
	}

	if err := p.Data.Machines().GetOrCreate(machine); err != nil {
		return nil, fmt.Errorf("failed to get or create machine: %v", err)
	}

	fmt.Printf("machine %#+v\n", machine)

	for i, _ := range scanReport.Files {
		scanReport.Files[i].ReportRunID = reportRun.ID
		scanReport.Files[i].MachineID = machine.ID
	}

	if err := p.Data.ScannedFiles().InsertBatch(scanReport.Files); err != nil {

		return nil, fmt.Errorf("failed to insert scanned files: %v", err)
	}

	return reportRun, nil
}

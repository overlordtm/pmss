package pmss

import (
	"fmt"

	"github.com/overlordtm/pmss/internal/utils"
	"github.com/overlordtm/pmss/pkg/datastore"
	"gorm.io/gorm"
)

type Pmss struct {
	dbUrl string

	db *gorm.DB
}

type Option func(*Pmss)

type HashSearchResult struct {
	File *datastore.KnownFile
}

func WithDbUrl(dbUrl string) Option {
	return func(p *Pmss) {
		p.dbUrl = dbUrl
	}
}

func New(options ...Option) (*Pmss, error) {

	pmms := &Pmss{}

	for _, option := range options {
		option(pmms)
	}

	dialector, err := utils.ParseDBUrl(pmms.dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database url: %v", err)
	}

	if db, err := gorm.Open(dialector, &gorm.Config{}); err != nil {
		return nil, fmt.Errorf("failed to initialize datastore: %v", err)
	} else {
		pmms.db = db
	}

	datastore.AutoMigrate(pmms.db)

	return pmms, nil
}

func (p *Pmss) FindByHash(hash string) (*HashSearchResult, error) {
	knownFile := new(datastore.KnownFile)

	if err := datastore.KnownFiles().FindByHash(hash, knownFile)(p.db); err != nil {
		return nil, err
	}
	return &HashSearchResult{
		File: knownFile,
	}, nil
}

func (p *Pmss) FindMachineByHostname(machineHostname string, machine *datastore.Machine) (bool, error) {
	if err := datastore.Machines().FindByHostname(machineHostname, machine)(p.db); err != nil {
		return false, err
	}
	return true, nil
}

func (p *Pmss) DoMachineReport(scanReport *ScanReport) (*datastore.ReportRun, error) {

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var reportRun *datastore.ReportRun

	if scanReport.ScanRunId != nil {
		if err := datastore.ReportRuns().FindByID(*scanReport.ScanRunId, reportRun)(tx); err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return nil, fmt.Errorf("failed to find report run: %v", err)
		}
	}

	if reportRun == nil {
		reportRun = &datastore.ReportRun{}
		if err := datastore.ReportRuns().Create(reportRun)(tx); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create new report run: %v", err)
		}
	}

	var machine *datastore.Machine = &datastore.Machine{
		MachineId: scanReport.MachineId,
	}

	if err := datastore.Machines().FirstOrCreate(machine)(tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create machine: %v", err)
	}

	for i, _ := range scanReport.Files {
		scanReport.Files[i].ReportRunID = reportRun.ID
		scanReport.Files[i].MachineID = machine.ID
	}

	if err := datastore.ScannedFiles().CreateInBatches(scanReport.Files)(tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert scanned files: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return reportRun, nil
}

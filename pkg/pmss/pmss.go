package pmss

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/overlordtm/pmss/pkg/apitypes"
	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/pkgscraper"
	"gorm.io/gorm"
)

type Pmss struct {
	dbUrl string

	db *gorm.DB
}

type Option func(*Pmss)

type HashSearchResult apitypes.KnownFile

func WithDbUrl(dbUrl string) Option {
	return func(p *Pmss) {
		p.dbUrl = dbUrl
	}
}

func New(options ...Option) (*Pmss, error) {

	pmss := &Pmss{}

	for _, option := range options {
		option(pmss)
	}

	if db, err := datastore.Open(pmss.dbUrl); err != nil {
		return nil, fmt.Errorf("failed to initialize datastore: %v", err)
	} else {
		pmss.db = db
	}

	datastore.AutoMigrate(pmss.db)

	return pmss, nil
}

func (p *Pmss) FindByHash(hash string) (*HashSearchResult, error) {
	knownFile := new(datastore.KnownFile)

	if err := datastore.KnownFiles().FindByHash(hash, knownFile)(p.db); err != nil {
		return nil, err
	}
	return &HashSearchResult{
		KnownPath: knownFile.Path,
		Status:    knownFile.Status,
	}, nil
}

func (p *Pmss) FindByHashBatch(hash []apitypes.HashQuery) ([]HashSearchResult, error) {

	var knownFiles []HashSearchResult

	for _, item := range hash {
		knownFile := new(datastore.KnownFile)

		pth := item.Path

		if err := datastore.KnownFiles().FindByHash(item.Hash, knownFile)(p.db); err != nil {
			switch err {
			case gorm.ErrRecordNotFound:

				knownFiles = append(knownFiles, HashSearchResult{Path: pth, Status: datastore.FileStatusUnknown})
			default:
				return nil, err
			}

		} else {
			knownFiles = append(knownFiles, HashSearchResult{
				Path:      pth,
				Size:      knownFile.Size,
				KnownPath: knownFile.Path,
				Status:    knownFile.Status,
				Md5:       knownFile.MD5,
				Sha1:      knownFile.SHA1,
				Sha256:    knownFile.SHA256,
			})
		}
	}

	return knownFiles, nil
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

	var machine *datastore.Machine = &datastore.Machine{
		MachineId: scanReport.MachineId,
		Hostname:  scanReport.Hostname,
	}

	if err := datastore.Machines().FirstOrCreate(machine)(tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create machine: %v", err)
	}

	var reportRun *datastore.ReportRun = &datastore.ReportRun{
		MachineID: machine.ID,
		IP:        scanReport.IP,
	}

	if scanReport.ScanRunId != nil {
		reportRun.ID = *scanReport.ScanRunId
	}

	if err := datastore.ReportRuns().FirstOrCreate(reportRun)(tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create report run: %v", err)
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

func (p *Pmss) ScanFile(scannedFile *datastore.ScannedFile) (*datastore.KnownFile, error) {
	knownFile := &datastore.KnownFile{
		Status: datastore.FileStatusUnknown,
	}
	err := datastore.KnownFiles().FindByScannedFile(scannedFile, knownFile)(p.db)
	if err != nil {
		return knownFile, err
	}
	return knownFile, nil
}

func (p *Pmss) UpdatePackages(ctx context.Context) (err error) {

	releases := []string{"bullseye", "bullseye-updates", "bullseye-security", "bullseye-backports"}
	components := []string{"main", "non-free", "contrib"}
	for _, release := range releases {
		for _, component := range components {
			if err1 := p.db.Transaction(func(tx *gorm.DB) error {
				return pkgscraper.ScrapeDebianMirror(ctx, p.db, release, "amd64", component)
			}); err1 != nil {
				err = errors.Join(err, err1)
			}
		}
	}

	releases = []string{"focal", "focal-updates", "focal-security", "focal-backports", "jammy", "jammy-updates", "jammy-security", "jammy-backports"}
	components = []string{"main", "multiverse", "restricted", "universe"}

	for _, release := range releases {
		for _, component := range components {

			if err1 := p.db.Transaction(func(tx *gorm.DB) error {
				return pkgscraper.ScrapeUbuntuMirror(ctx, p.db, release, "amd64", component)
			}); err1 != nil {
				err = errors.Join(err, err1)
			}
		}
	}

	return err
}

func (p *Pmss) UpdatePackageHashes(ctx context.Context, concurrency int) error {

	packages := make([]datastore.Package, 0)

	if err := p.db.Model(&datastore.Package{}).Where("scraped_at IS NULL").Find(&packages).Error; err != nil {
		return fmt.Errorf("failed to get unscraped packages: %v", err)
	}

	workCh := make(chan datastore.Package, 1024)

	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	go func() {
		for _, pkg := range packages {
			workCh <- pkg
		}
		close(workCh)
	}()

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for pkg := range workCh {
				p.db.Transaction(func(tx *gorm.DB) error {
					if err := pkgscraper.ScrapeDebianPackage(ctx, p.db, pkg); err != nil {
						return fmt.Errorf("failed to scrape package: %w", err)
					}

					pkg.ScrapedAt = time.Now()

					return datastore.Packages().Save(pkg)(p.db)
				})
			}
		}()
	}

	wg.Wait()

	return nil
}

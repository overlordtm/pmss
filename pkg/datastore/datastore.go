package datastore

import (
	"fmt"

	"gorm.io/gorm"
)

// Store is a datastore
type Store struct {
	opts options

	packageRepository     *packageRepository
	machineRepository     *machineRepository
	knownFileRepository   *knownFileRepository
	scannedFileRepository *scannedFileRepository
	reportRunRepository   *reportRunRepository
}

type options struct {
	dialector gorm.Dialector
}

type Option func(*options)

func WithDb(dialector gorm.Dialector) Option {
	return func(o *options) {
		o.dialector = dialector
	}
}

func New(opts ...Option) (*Store, error) {

	o := options{}

	for _, option := range opts {
		option(&o)
	}

	db, err := gorm.Open(o.dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %v", err)
	}

	db.AutoMigrate(&Package{})
	db.AutoMigrate(&Machine{})
	db.AutoMigrate(&KnownFile{})
	db.AutoMigrate(&ScannedFile{})
	db.AutoMigrate(&ReportRun{})

	return &Store{
		opts:                  o,
		packageRepository:     &packageRepository{db},
		machineRepository:     &machineRepository{db},
		knownFileRepository:   &knownFileRepository{db},
		scannedFileRepository: &scannedFileRepository{db},
		reportRunRepository:   &reportRunRepository{db},
	}, nil
}

func (ds *Store) Packages() *packageRepository {
	return ds.packageRepository
}

func (ds *Store) Machines() *machineRepository {
	return ds.machineRepository
}

func (ds *Store) KnownFiles() *knownFileRepository {
	return ds.knownFileRepository
}

func (ds *Store) ScannedFiles() *scannedFileRepository {
	return ds.scannedFileRepository
}

func (ds *Store) ReportRuns() *reportRunRepository {
	return ds.reportRunRepository
}

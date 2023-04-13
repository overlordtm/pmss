package datastore

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/overlordtm/pmss/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Datastore interface {
	Machines() *machineRepository
	Packages() *packageRepository
	KnownFiles() *knownFileRepository
	ScannedFiles() *scannedFileRepository
	ReportRuns() *reportRunRepository
}

type DbOp func(*gorm.DB) error

func Open(dbUrl string) (*gorm.DB, error) {
	dialector, err := utils.ParseDBUrl(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database url: %v", err)
	}

	if db, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             1000 * time.Millisecond,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
				LogLevel:                  logger.Warn,
			},
		),
	}); err != nil {
		return nil, fmt.Errorf("failed to initialize datastore: %v", err)
	} else {
		return db, nil
	}
}

func MustOpen(dbUrl string) *gorm.DB {
	db, err := Open(dbUrl)
	if err != nil {
		panic(err)
	}
	return db
}

func AutoMigrate(db *gorm.DB) (err error) {

	err = errors.Join(err, db.AutoMigrate(&Machine{}))
	err = errors.Join(err, db.AutoMigrate(&Package{}))
	err = errors.Join(err, db.AutoMigrate(&KnownFile{}))
	err = errors.Join(err, db.AutoMigrate(&ScannedFile{}))
	err = errors.Join(err, db.AutoMigrate(&ReportRun{}))

	return err
}

func MustAutoMigrate(db *gorm.DB) {
	if err := AutoMigrate(db); err != nil {
		panic(err)
	}
}

func Packages() *packageRepository {
	return &packageRepository{}
}

func Machines() *machineRepository {
	return &machineRepository{}
}

func KnownFiles() *knownFileRepository {
	return &knownFileRepository{}
}

func ScannedFiles() *scannedFileRepository {
	return &scannedFileRepository{}
}

func ReportRuns() *reportRunRepository {
	return &reportRunRepository{}
}

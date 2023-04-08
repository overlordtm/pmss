package datastore

import (
	"gorm.io/gorm"
)

// ReportRun represents information about a machine on the network. It also contains info whether the machine is allowed to submit files.
type ReportRun struct {
	gorm.Model
	Files []ScannedFile `gorm:"foreignKey:ReportRunID"`
}

type reportRunRepository struct {
	db *gorm.DB
}

func (repo *machineRepository) NewReportRun() (*ReportRun, error) {
	reportRun := &ReportRun{}
	if err := repo.db.Create(reportRun).Error; err != nil {
		return nil, err
	}
	return reportRun, nil
}
func (repo *reportRunRepository) Insert(row Machine) error {
	return repo.db.Create(&row).Error
}

func (repo *reportRunRepository) InsertBatch(rows []Machine) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

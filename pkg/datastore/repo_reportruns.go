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

func (repo *reportRunRepository) CreateNew(dest *ReportRun) error {
	*dest = ReportRun{}
	return repo.db.Create(dest).Error
}
func (repo *reportRunRepository) Insert(row Machine) error {
	return repo.db.Create(&row).Error
}

func (repo *reportRunRepository) InsertBatch(rows []Machine) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

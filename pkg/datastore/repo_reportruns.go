package datastore

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRun represents information about a machine on the network. It also contains info whether the machine is allowed to submit files.
type ReportRun struct {
	ID        uuid.UUID `gorm:"type:char(36);primarykey;default:uuid()"`
	CreatedAt time.Time
	Files     []ScannedFile `gorm:"foreignKey:ReportRunID"`
}

type reportRunRepository struct {
	db *gorm.DB
}

func (repo *reportRunRepository) FindByID(id uuid.UUID, dest *ReportRun) error {
	return repo.db.First(dest, id).Error
}

func (r *reportRunRepository) GetOrCreate(reportRun *ReportRun) error {
	if err := r.db.FirstOrCreate(reportRun, "id = ?", reportRun.ID).Error; err != nil {
		return err
	}
	return nil
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

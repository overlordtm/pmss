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

type reportRunRepository struct{}

func (*reportRunRepository) FindByID(id uuid.UUID, dest *ReportRun) DbOp {
	return func(d *gorm.DB) error {
		return d.Where("id = ?", id).First(dest).Error
	}
}

func (*reportRunRepository) FirstOrCreate(reportRun *ReportRun) DbOp {
	return func(d *gorm.DB) error {
		return d.FirstOrCreate(reportRun, "id = ?", reportRun.ID).Error
	}
}

func (*reportRunRepository) Create(dest *ReportRun) DbOp {
	return func(d *gorm.DB) error {
		*dest = ReportRun{}
		return d.Create(dest).Error
	}
}

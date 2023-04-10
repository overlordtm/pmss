package datastore

import (
	"time"

	"github.com/google/uuid"
	"github.com/overlordtm/pmss/pkg/hashtype"
	"gorm.io/gorm"
)

type ScannedFile struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Path
	Path string `gorm:"type:varchar(1024);notnull"`

	// Hashes
	SHA1   *string `gorm:"type:char(40);check:at_least_one,(sha1 IS NOT NULL) OR (sha256 IS NOT NULL) OR (md5 IS NOT NULL)"`
	SHA256 *string `gorm:"type:char(64)"`
	MD5    *string `gorm:"type:char(32)"`

	// File info
	Size     int64  `gorm:"notnull;default:-1"`
	Mode     uint32 `gorm:"notnull;default:4294967295"`
	MimeType string `gorm:"type:varchar(255);notnull"`

	// File times
	Mtime *int64 `gorm:"type:long"`
	Atime *int64 `gorm:"type:long"`
	Ctime *int64 `gorm:"type:long"`

	// Users
	Owner string `gorm:"type:varchar(255);notnull"`
	Group string `gorm:"type:varchar(255);notnull"`

	// Known file reference
	KnownMatchID *uint
	KnownMatch   KnownFile `gorm:"foreignKey:KnownMatchID"`

	// Run info
	ReportRunID uuid.UUID
	ReportRun   ReportRun `gorm:"foreignKey:ReportRunID"`

	// Submitting Machine info
	MachineID uint
	Machine   Machine `gorm:"foreignKey:MachineID"`
}

type scannedFileRepository struct {
	db *gorm.DB
}

func (*scannedFileRepository) FindByHash(hash string, dest *KnownFile) DbOp {
	return func(d *gorm.DB) error {
		switch hashtype.Detect(hash) {
		case hashtype.SHA1:
			return d.Where("sha1 = ?", hash).First(dest).Error
		case hashtype.SHA256:
			return d.Where("sha256 = ?", hash).First(dest).Error
		case hashtype.MD5:
			return d.Where("md5 = ?", hash).First(dest).Error
		default:
			return hashtype.ErrUnknown
		}
	}
}

func (*scannedFileRepository) Create(row ScannedFile) DbOp {
	return func(d *gorm.DB) error {
		return d.Create(&row).Error
	}
}

func (*scannedFileRepository) CreateInBatches(rows []ScannedFile) DbOp {
	return func(d *gorm.DB) error {
		return d.CreateInBatches(&rows, 1000).Error
	}
}

package datastore

import (
	"github.com/google/uuid"
	"github.com/overlordtm/pmss/pkg/hashvariant"
	"gorm.io/gorm"
)

type ScannedFile struct {
	gorm.Model

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

func (repo *scannedFileRepository) prepFindByHash(hash string, variant *hashvariant.HashVariant) (*gorm.DB, error) {
	*variant = hashvariant.DetectHashVariant(hash)
	switch *variant {
	case hashvariant.SHA1:
		return repo.db.Where("sha1 = ?", hash), nil
	case hashvariant.SHA256:
		return repo.db.Where("sha256 = ?", hash), nil
	case hashvariant.MD5:
		return repo.db.Where("md5 = ?", hash), nil
	}
	//return repo.db.Where("1 = 0")
	return nil, hashvariant.ErrUnknownHashVariant
}

func (repo *scannedFileRepository) FindByHash(hash string, dest *KnownFile, destVariant *hashvariant.HashVariant) error {
	stmt, err := repo.prepFindByHash(hash, destVariant)
	if err != nil {
		return err
	}
	return stmt.First(dest).Error
}

func (repo *scannedFileRepository) FindAllByHash(hash string, dest *KnownFile, destVariant *hashvariant.HashVariant) error {
	stmt, err := repo.prepFindByHash(hash, destVariant)
	if err != nil {
		return err
	}
	return stmt.Find(dest).Error
}

func (repo *scannedFileRepository) DB() *gorm.DB {
	return repo.db
}

func (repo *scannedFileRepository) Insert(row ScannedFile) error {
	return repo.db.Create(&row).Error
}

func (repo *scannedFileRepository) InsertBatch(rows []ScannedFile) error {
	return repo.db.CreateInBatches(&rows, 1000).Error
}

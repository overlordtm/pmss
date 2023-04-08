package datastore

import (
	"gorm.io/gorm"
)

type ScannedFile struct {
	gorm.Model

	// Path
	Path string `gorm:"type:varchar(1024);notnull"`

	// Hashes
	SHA1   string `gorm:"type:char(40);check:at_least_one,(sha1 IS NOT NULL) OR (sha256 IS NOT NULL) OR (md5 IS NOT NULL)"`
	SHA256 string `gorm:"type:char(64)"`
	MD5    string `gorm:"type:char(32)"`

	// File info
	Size     uint64 `gorm:"notnull"`
	FileMode uint32 `gorm:"notnull"`
	MimeType string `gorm:"type:varchar(255);notnull"`

	// File times
	Mtime uint64 `gorm:"type:timestamp"`
	Atime uint64 `gorm:"type:timestamp"`
	Ctime uint64 `gorm:"type:timestamp"`

	// Users
	Owner string `gorm:"type:varchar(255);notnull"`
	Group string `gorm:"type:varchar(255);notnull"`

	// Known file reference
	KnownMatchID *uint
	KnownMatch   KnownFile `gorm:"foreignKey:KnownMatchID"`

	// Run info
	ReportRunID uint
	ReportRun   ReportRun `gorm:"foreignKey:ReportRunID"`

	// Submitting Machine info
	MachineID uint
	Machine   Machine `gorm:"foreignKey:MachineID"`
}

type scannedFileRepository struct {
	db *gorm.DB
}

func (repo *scannedFileRepository) FindBySHA1(sha1 string) (*ScannedFile, error) {
	return repo.FindBy(&ScannedFile{SHA1: sha1})
}

func (repo *scannedFileRepository) FindBy(fields *ScannedFile) (*ScannedFile, error) {
	var row ScannedFile
	if err := repo.db.Where(fields).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (repo *scannedFileRepository) Insert(row ScannedFile) error {
	return repo.db.Create(&row).Error
}

func (repo *scannedFileRepository) InsertBatch(rows []ScannedFile) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

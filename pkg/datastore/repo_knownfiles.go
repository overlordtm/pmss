package datastore

import (
	"gorm.io/gorm"
)

type KnownFile struct {
	gorm.Model
	// Path, Hashes, Indexed
	Path   string `gorm:"type:varchar(4096),index:path_sha1"`
	SHA1   string `gorm:"type:char(40),index:path_sha1"`
	SHA256 string `gorm:"type:char(64)"`
	MD5    string `gorm:"type:char(32)"`

	// File info
	Size     int64  `gorm:"type:bigint"`
	Mode     uint32 `gorm:"type:unsigned int"`
	MimeType string `gorm:"type:varchar(255)"`

	// File times
	Mtime uint64 `gorm:"type:timestamp"`
	Atime uint64 `gorm:"type:timestamp"`
	Ctime uint64 `gorm:"type:timestamp"`

	// Wether was scraped or voted for
	FromDeb bool `gorm:"type:bool"`

	// Is it malicious
	IsSafe bool `gorm:"type:bool"`
}

type knownFileRepository struct {
	db *gorm.DB
}

func (repo *knownFileRepository) FindByMD5(md5 string) (*KnownFile, error) {
	var row KnownFile
	err := repo.db.Where("md5 = ?", md5).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (repo *knownFileRepository) Insert(row KnownFile) error {
	return repo.db.Create(&row).Error
}

func (repo *knownFileRepository) InsertBatch(rows []KnownFile) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

func (repo *knownFileRepository) List() ([]KnownFile, error) {
	var rows []KnownFile
	if err := repo.db.Find(rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

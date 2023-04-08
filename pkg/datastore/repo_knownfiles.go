package datastore

import (
	"github.com/overlordtm/pmss/pkg/hashvariant"
	"gorm.io/gorm"
)

type KnownFile struct {
	gorm.Model
	// Path, Hashes, Indexed
	Path   string `gorm:"varchar(1024);index:path_sha1;notnull"`
	SHA1   string `gorm:"type:char(40);index:path_sha1;notnull"`
	SHA256 string `gorm:"type:char(64);notnull"`
	MD5    string `gorm:"type:char(32);notnull"`

	// File info
	Size     int64  `gorm:"notnull;default:null"`
	Mode     uint32 `gorm:"notnull;default:null"`
	MimeType string `gorm:"type:char(255);notnull;default:null"`

	// File times
	Mtime uint64 `gorm:"type:timestamp"`
	Atime uint64 `gorm:"type:timestamp"`
	Ctime uint64 `gorm:"type:timestamp"`

	// Wether was scraped or voted for
	FromDeb bool `gorm:"notnull;default:false"`

	// Is it malicious
	IsSafe bool `gorm:"notnull"`
}

type knownFileRepository struct {
	db *gorm.DB
}

func (repo *knownFileRepository) prepFindByHash(hash string, destVariant *hashvariant.HashVariant) (ctx *gorm.DB) {
	*destVariant = hashvariant.DetectHashVariant(hash)
	switch *destVariant {
	case hashvariant.SHA1:
		ctx = repo.db.Where("sha1 = ?", hash)
	case hashvariant.SHA256:
		ctx = repo.db.Where("sha256 = ?", hash)
	case hashvariant.MD5:
		ctx = repo.db.Where("md5 = ?", hash)
	default:
		ctx = repo.db.Where("FALSE")
	}
	return
}

func (repo *knownFileRepository) FindByHash(hash string, dest *KnownFile, destVariant *hashvariant.HashVariant) error {
	return repo.prepFindByHash(hash, destVariant).First(dest).Error
}

func (repo *knownFileRepository) FindAllPathsByHash(hash string, dest *[]string, destVariant *hashvariant.HashVariant) error {
	return repo.prepFindByHash(hash, destVariant).Model(&ScannedFile{}).Select("path").Find(dest).Error
}

func (repo *knownFileRepository) FindAllByHash(hash string, dest *[]KnownFile, destVariant *hashvariant.HashVariant) error {
	return repo.prepFindByHash(hash, destVariant).Find(dest).Error
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

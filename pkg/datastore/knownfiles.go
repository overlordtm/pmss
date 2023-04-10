package datastore

import (
	"fmt"

	"github.com/overlordtm/pmss/pkg/hashtype"
	"gorm.io/gorm"
)

type FileStatus byte

const (
	FileStatusMalicious FileStatus = 1 << iota
	FileStatusSafe
	FileStatusUnknown
)

type KnownFile struct {
	*gorm.Model
	// Path, Hashes, Indexed
	// Path   *string `gorm:"varchar(1024);index:path"`
	SHA1   *string `gorm:"type:char(40)"`
	SHA256 *string `gorm:"type:char(64)"`
	MD5    *string `gorm:"type:char(32)"`

	// File info
	Size     *int64  `gorm:"default:null"`
	MimeType *string `gorm:"default:null"`

	// Wether was scraped or voted for
	FromDeb bool `gorm:"notnull;default:false"`

	// File Status
	Status FileStatus `gorm:"notnull;default:1"`
}

func (f *KnownFile) BeforeCreate(tx *gorm.DB) (err error) {
	if f.SHA1 == nil && f.SHA256 == nil && f.MD5 == nil {
		err = fmt.Errorf("at least one hash must be provided")
		tx.AddError(err)
		return err
	}
	return nil
}

type knownFileRepository struct {
}

func (*knownFileRepository) FindByHash(hash string, dest *KnownFile) DbOp {
	return func(d *gorm.DB) error {
		hashType := hashtype.Detect(hash)
		switch hashType {
		case hashtype.SHA1:
			return d.Model(&KnownFile{}).Where("sha1 = ?", hash).First(dest).Error
		case hashtype.SHA256:
			return d.Model(&KnownFile{}).Where("sha256 = ?", hash).First(dest).Error
		case hashtype.MD5:
			return d.Model(&KnownFile{}).Where("md5 = ?", hash).First(dest).Error
		default:
			fmt.Println(fmt.Errorf("detect hash type: %v %s", hashType, hash))
			err := hashtype.ErrUnknown
			d.AddError(err)
			return err
		}
	}
}

func (*knownFileRepository) Create(row KnownFile) DbOp {
	return func(db *gorm.DB) error {
		return db.Create(&row).Error
	}
}

func (*knownFileRepository) CreateInBatches(rows []KnownFile) DbOp {
	return func(db *gorm.DB) error {
		return db.CreateInBatches(&rows, 100).Error
	}
}

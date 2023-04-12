// go:generate go run internal/gen/fixturegen/main.go

package datastore

import (
	"fmt"

	"github.com/overlordtm/pmss/pkg/hashtype"
	"gorm.io/gorm"
)

type FileStatus byte

const (
	FileStatusUnknown   FileStatus = 0
	FileStatusSafe      FileStatus = 1
	FileStatusMalicious FileStatus = 255
)

func (f FileStatus) String() string {
	switch f {
	case FileStatusUnknown:
		return "unknown"
	case FileStatusSafe:
		return "safe"
	case FileStatusMalicious:
		return "malicious"
	default:
		return "unknown"
	}
}

type KnownFile struct {
	ID uint `gorm:"primarykey"`

	// Path, Hashes, Indexed
	Path *string `gorm:"varchar(1024)"`

	SHA1   *string `gorm:"type:char(40);index"`
	SHA256 *string `gorm:"type:char(64);index"`
	MD5    *string `gorm:"type:char(32);index"`

	// File info
	Size     *int64  `gorm:"default:null"`
	MimeType *string `gorm:"default:null"`

	// Wether was scraped or voted for
	PackageID *uint
	Package   Package `gorm:"foreignKey:PackageID"`

	// File Status
	Status FileStatus `gorm:"notnull;default:0"`
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

func (*knownFileRepository) All(all []KnownFile) DbOp {
	return func(d *gorm.DB) error {
		return d.Model(&KnownFile{}).Find(&all).Error
	}
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

func (*knownFileRepository) FindByScannedFile(scannedFile *ScannedFile, dest *KnownFile) DbOp {
	return func(d *gorm.DB) error {
		q := d.Model(&KnownFile{})

		if scannedFile.SHA1 != nil {
			q = q.Or("sha1 = ?", scannedFile.SHA1)
		}
		if scannedFile.SHA256 != nil {
			q = q.Or("sha256 = ?", scannedFile.SHA256)
		}
		if scannedFile.MD5 != nil {
			q = q.Or("md5 = ?", scannedFile.MD5)
		}
		q.Statement.RaiseErrorOnNotFound = true
		return q.Limit(1).Order("status DESC").Find(dest).Error
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

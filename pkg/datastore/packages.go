package datastore

import (
	"time"

	"gorm.io/gorm"
)

type OsType string

const (
	OsTypeDebian    OsType = "debian"
	OsTypeRedhat    OsType = "redhat"
	OsTypeUbuntu    OsType = "ubuntu"
	OsTypeAlmaLinux OsType = "almalinux"
)

type Package struct {
	ID uint `gorm:"primarykey"`

	Name     string `gorm:"uniqueIndex:package_name_version_architecture"`
	Filename string
	MD5Sum   *string `gorm:"default:null"`
	SHA256   *string `gorm:"default:null"`
	Version  string  `gorm:"uniqueIndex:package_name_version_architecture"`
	Size     uint64

	Architecture string `gorm:"uniqueIndex:package_name_version_architecture"`
	Distro       string
	Component    string
	OsType       OsType

	CreatedAt time.Time
	UpdatedAt time.Time
	ScrapedAt time.Time `gorm:"default:null"`
}

type packageRepository struct {
}

func (*packageRepository) FindAll(dest []Package) DbOp {
	return func(d *gorm.DB) error {
		return d.Model(&Package{}).Find(&dest).Error
	}
}

func (*packageRepository) Save(p Package) DbOp {
	return func(d *gorm.DB) error {
		return d.Save(&p).Error
	}
}

func (*packageRepository) CreateInBatches(packages []Package) DbOp {
	return func(d *gorm.DB) error {
		return d.CreateInBatches(packages, 1000).Error
	}
}

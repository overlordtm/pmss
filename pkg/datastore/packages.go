package datastore

import (
	"time"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Name         string
	Version      string
	Size         uint64
	MD5          string
	SHA256       string
	Architecture string
	Filename     string
	Distro       string
	ScrapedAt    *time.Time
}

type packageRepository struct {
	db *gorm.DB
}

func (r *packageRepository) FindByMD5(md5 string) (*Package, error) {
	var row Package
	err := r.db.Where("md5 = ?", md5).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *packageRepository) FindByName(name string) (*Package, error) {
	var row Package
	err := r.db.Where("name = ?", name).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *packageRepository) FindAll() ([]Package, error) {
	var rows []Package
	err := r.db.Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *packageRepository) Insert(row Package) error {
	return r.db.Create(&row).Error
}

func (r *packageRepository) InsertBatch(rows []Package) error {
	return r.db.CreateInBatches(&rows, 100).Error
}

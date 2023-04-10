package datastore

import (
	"gorm.io/gorm"
)

type Package struct {
	Name         string
	Version      string
	Size         uint64
	MD5          string
	SHA256       string
	Architecture string
	Filename     string
	Distro       string
	// ScrapedAt    *time.Time
}

type packageRepository struct {
	db *gorm.DB
}

func (repo *packageRepository) FindByMD5(md5 string) (*Package, error) {
	var row Package
	err := repo.db.Where("md5 = ?", md5).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (repo *packageRepository) FindByName(name string) (*Package, error) {
	var row Package
	err := repo.db.Where("name = ?", name).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (repo *packageRepository) FindAll() ([]Package, error) {
	var rows []Package
	err := repo.db.Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *packageRepository) Insert(row Package) error {
	return repo.db.Create(&row).Error
}

func (repo *packageRepository) InsertBatch(rows []Package) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

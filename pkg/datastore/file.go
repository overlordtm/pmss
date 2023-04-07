package datastore

import (
	"gorm.io/gorm"
)

type File struct {
	*gorm.Model
	MachineID string `gorm:"type:varchar(255)"`
	MachineIP string `gorm:"type:varchar(255)"`
	Path      string `gorm:"type:varchar(4096)"`
	MD5       string `gorm:"type:char(32)"`
	SHA1      string `gorm:"type:char(40)"`
	SHA256    string `gorm:"type:char(64)"`
	Size      uint64
	Mtime     uint64
	Atime     uint64
	Ctime     uint64
	FileMode  uint32
	MimeType  string `gorm:"type:varchar(255)"`
	// IsSafe    bool
}

type fileRepository struct {
	db *gorm.DB
}

func (repo *fileRepository) FindByMD5(md5 string) (*File, error) {
	var row File
	err := repo.db.Where("md5 = ?", md5).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

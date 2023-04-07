package datastore

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type WhitelistMeta struct {
	Package string `json:"package"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
	Distro  string `json:"distro"`
	Size    int64  `json:"size"`
	Mode    uint32 `json:"mode"`
	Owner   string `json:"owner"`
	Group   string `json:"group"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *WhitelistMeta) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, j)
}

// Value return json value, implement driver.Valuer interface
func (j WhitelistMeta) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (WhitelistMeta) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

// WhitelistItem represents a file that is know to be good (non-malicious). Less important attributes (the ones you do not need to query by) are stored in Meta field as JSON
type WhitelistItem struct {
	gorm.Model
	MD5    string
	SHA1   string
	SHA256 string
	Path   string
	Meta   WhitelistMeta
}

type whitelistRepository struct {
	db *gorm.DB
}

func (repo *whitelistRepository) FindByMD5(md5 string) (*WhitelistItem, error) {
	var row WhitelistItem
	err := repo.db.Where("md5 = ?", md5).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (repo *whitelistRepository) Insert(row WhitelistItem) error {
	return repo.db.Create(&row).Error
}

func (repo *whitelistRepository) InsertBatch(rows []WhitelistItem) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

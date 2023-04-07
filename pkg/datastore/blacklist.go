package datastore

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type BlacklistMeta struct {
	Platform string `json:"platform"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *BlacklistMeta) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, j)
}

// Value return json value, implement driver.Valuer interface
func (j BlacklistMeta) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (BlacklistMeta) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

type BlacklistItem struct {
	gorm.Model
	MD5       string
	SHA1      string
	SHA256    string
	Signature string
	Meta      BlacklistMeta
}

type blacklistRepository struct {
	db *gorm.DB
}

func (repo *blacklistRepository) FindByMD5(md5 string) (*BlacklistItem, error) {
	var row BlacklistItem
	err := repo.db.Where("md5 = ?", md5).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

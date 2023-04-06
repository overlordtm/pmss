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

type WhitelistItem struct {
	gorm.Model
	MD5    string
	SHA1   string
	SHA256 string
	Path   string
	Meta   WhitelistMeta
}

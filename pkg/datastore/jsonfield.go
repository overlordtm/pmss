package datastore

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSONField[T any] struct {
	Val T
}

func (j JSONField[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Val)
}

func (j *JSONField[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &j.Val)
}

func (j JSONField[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Val)
}

func (j *JSONField[T]) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		return errors.New("JSONField: type assertion to []byte failed")
	}

	return json.Unmarshal(source, &j.Val)
}

func (JSONField[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

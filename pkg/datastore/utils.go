package datastore

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ParseDBUrl(dbUrl string) (dialector gorm.Dialector, err error) {

	if strings.HasPrefix(dbUrl, "sqlite3://") {
		return sqlite.Open(dbUrl[len("sqlite3://"):]), nil
	} else if strings.HasPrefix(dbUrl, "mysql://") {
		return mysql.Open(dbUrl[len("mysql://"):]), nil
	} else {
		return nil, fmt.Errorf("unsupported database scheme: %s", dbUrl)
	}
}

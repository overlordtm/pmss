package testutils

import (
	"github.com/overlordtm/pmss/pkg/datastore"
	"gorm.io/gorm"
)

func MustExecute(op datastore.DbOp, db *gorm.DB) {
	if err := op(db); err != nil {
		panic(err)
	}
}

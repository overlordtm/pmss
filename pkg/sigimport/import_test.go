package sigimport

import (
	"fmt"
	"os"
	"testing"
)

func TestImport(t *testing.T) {

	dbPath := "test.db"
	defer os.Remove(dbPath)

	err := Import(fmt.Sprintf("%s?_journal=OFF&_locking=EXCLUSIVE&_sync=OFF", dbPath), "/home/az/ws.az/pmss/data/full2.csv")
	if err != nil {
		t.Error(err)
	}
}

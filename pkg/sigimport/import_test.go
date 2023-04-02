package sigimport

import (
	"fmt"
	"os"
	"testing"
)

func TestImport(t *testing.T) {
	projRoot := os.Getenv("PMSS_PROJ_ROOT")
	dbPath := fmt.Sprintf("%s/test/test.db", projRoot)
	fmt.Printf("dbPath: '%s'\n", dbPath)
	defer os.Remove(dbPath)

	csvFile := fmt.Sprintf("%s/test/test.csv", projRoot)
	fmt.Printf("csvFile: '%s'\n", csvFile)

	dbURI := fmt.Sprintf("%s?_journal=OFF&_locking=EXCLUSIVE&_sync=OFF", dbPath)
	err := Import(dbURI, csvFile)
	if err != nil {
		t.Error(err)
	}
}

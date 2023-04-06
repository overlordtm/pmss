package scanner

import (
	"fmt"
	"os"
	"testing"

	"github.com/overlordtm/pmss/pkg/checker/sigchecker"
	"github.com/overlordtm/pmss/pkg/sigdb"
)

func TestScan(t *testing.T) {
	// db := sigdb.New("/home/az/ws.az/pmss/test.db")
	// err := db.Init()
	// if err != nil {
	// 	t.Error(err)
	// }

	// c := sigchecker.New(db)

	// checkers := []interface{}{c}

	// scanner := New(db, checkers)

	// results, err := scanner.Scan("/home/az/ws.az/pmss/test/data")
	// if err != nil {
	// 	t.Error(err)
	// }

	// for _, r := range results {
	// 	if r.Err != nil {
	// 		t.Error(r.Err)
	// 	}
	// }
	t.Fatal("not implemented")
}

func TestScanFile(t *testing.T) {
	projRoot := os.Getenv("PMSS_PROJ_ROOT")
	dbPath := fmt.Sprintf("%s/test/0ad.db", projRoot)
	fmt.Printf("dbPath: '%s'\n", dbPath)

	db := sigdb.New(dbPath)
	err := db.Init()
	if err != nil {
		t.Error(err)
	}

	c := sigchecker.New(db)

	checkers := []interface{}{c}

	testFile := "0ad.appdata.xml"
	testFilePath := fmt.Sprintf("%s/test/data/%s", projRoot, testFile)

	r, err := scanFile(db, testFilePath, checkers)

	if err != nil {
		t.Error(err)
	}

	if r.Path != testFilePath {
		t.Error("wrong path")
	}

	for _, checkResult := range r.CheckResults {
		if checkResult.Signature != "Formbook" {
			t.Error("wrong signature")
		}
	}

	if len(r.CheckResults) != 1 {
		t.Error("wrong number of check results")
	}

}

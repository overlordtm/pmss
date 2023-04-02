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
	dbPath := fmt.Sprintf("%s/test/test.db", projRoot)
	fmt.Printf("dbPath: '%s'\n", dbPath)

	db := sigdb.New(dbPath)
	err := db.Init()
	if err != nil {
		t.Error(err)
	}

	c := sigchecker.New(db)

	checkers := []interface{}{c}

	pth := fmt.Sprintf("%s/test/data/179b98e2cb16a094755f853ae892b47948a8b6a83e7ca050d520e113ff180b2f.exe", projRoot)

	r, err := scanFile(db, pth, checkers)

	if err != nil {
		t.Error(err)
	}

	if r.Path != pth {
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

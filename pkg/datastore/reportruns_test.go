package datastore_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/overlordtm/pmss/pkg/datastore"
	"gorm.io/gorm"
)

func setup() *gorm.DB {
	db := datastore.MustOpen(testDbUrl)
	// drop all tables
	db.Migrator().DropTable(&datastore.ReportRun{})

	datastore.MustAutoMigrate(db)
	return db
}

func TestReportRuns(t *testing.T) {

	db := setup()

	// Create
	reportRun := datastore.ReportRun{
		ID: uuid.New(),
	}

	if err := datastore.ReportRuns().FirstOrCreate(&reportRun)(db); err != nil {
		t.Fatalf("failed to create report run: %v", err)
	}

	var count int64
	if err := db.Model(&reportRun).Count(&count).Error; err != nil {
		t.Fatal("failed to count report runs")
	} else {
		if count != 1 {
			t.Fatal("invalid count of report runs")
		}
	}

	if err := datastore.ReportRuns().FirstOrCreate(&reportRun)(db); err != nil {
		t.Fatalf("failed to create report run: %v", err)
	}

	if err := db.Model(&reportRun).Count(&count).Error; err != nil {
		t.Fatal("failed to count report runs")
	} else {
		if count != 1 {
			t.Fatal("invalid count of report runs")
		}
	}

}

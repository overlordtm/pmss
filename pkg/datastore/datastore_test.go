package datastore_test

import (
	"testing"

	"github.com/overlordtm/pmss/internal/utils"
	"github.com/overlordtm/pmss/pkg/datastore"
	"gorm.io/gorm"
)

func TestDatastore(t *testing.T) {
	testCases := []struct {
		testName string
		uri      string
	}{
		{
			"mysql",
			"mysql://pmss:pmss@tcp(mariadb:3306)/pmss?parseTime=true",
		},
		// {
		// 	"sqlite",
		// 	"sqlite3://:memory:",
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {

			dialector, err := utils.ParseDBUrl(testCase.uri)
			if err != nil {
				t.Error(err)
				return
			}

			db, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("error while opening database: %v", err)
			}

			datastore.AutoMigrate(db)

			rows := []datastore.Machine{
				{
					Hostname:  "hostname1.com",
					MachineId: "machineid1",
					IPv4:      utils.StringPtr("192.168.1.1"),
					IPv6:      utils.StringPtr("2001:fe8::1"),
				},
				{
					Hostname:  "hostname2.com",
					MachineId: "machineid2",
					IPv4:      utils.StringPtr("192.168.1.2"),
					IPv6:      utils.StringPtr("2001:fe8::2"),
				},
				{
					Hostname:  "hostname3.com",
					MachineId: "machineid3",
					IPv4:      utils.StringPtr("192.168.1.3"),
					IPv6:      utils.StringPtr("2001:fe8::3"),
				},
			}

			if err := datastore.Machines().CreateInBatches(rows)(db); err != nil {
				t.Error(err)
				return
			}
			var machine datastore.Machine
			if err := datastore.Machines().FindByIPv4("192.168.1.2", &machine)(db); err != nil {
				t.Error(err)
				return
			}
			var reportRun datastore.ReportRun
			if err := datastore.ReportRuns().Create(&reportRun)(db); err != nil {
				t.Error(err)
				return
			}
			scannedFiles := []datastore.ScannedFile{
				{
					Path:      "/bin/ps",
					SHA1:      utils.StringPtr("3c738552873525fda24139e1214c95bbdaf9dcca"),
					Machine:   machine,
					ReportRun: reportRun,
					Size:      137680,
					Mode:      0o755,
					Owner:     "root",
					Group:     "root",
				},
			}
			if err := datastore.ScannedFiles().CreateInBatches(scannedFiles)(db); err != nil {
				t.Error(err)
				return
			}

		})
	}
}

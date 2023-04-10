package datastore_test

import (
	"os"
	"testing"

	"github.com/overlordtm/pmss/internal/utils"
	"github.com/overlordtm/pmss/pkg/datastore"
)

var (
	testDbUrl string
)

func TestMain(m *testing.M) {
	testDbUrl = utils.EnvOrDefault("PMSS_TEST_DB_URL", "mysql://pmss:pmss@tcp(mariadb:3306)/test_pmss?parseTime=true")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestDatastore(t *testing.T) {
	testCases := []struct {
		testName string
		uri      string
	}{
		{
			"mysql",
			testDbUrl,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			tx := datastore.MustOpen(testCase.uri).Begin()
			datastore.AutoMigrate(tx)
			defer tx.Rollback()

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

			if err := datastore.Machines().CreateInBatches(rows)(tx); err != nil {
				t.Error(err)
				return
			}
			var machine datastore.Machine
			if err := datastore.Machines().FindByIPv4("192.168.1.2", &machine)(tx); err != nil {
				t.Error(err)
				return
			}
			var reportRun datastore.ReportRun
			if err := datastore.ReportRuns().Create(&reportRun)(tx); err != nil {
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
			if err := datastore.ScannedFiles().CreateInBatches(scannedFiles)(tx); err != nil {
				t.Error(err)
				return
			}

		})
	}
}

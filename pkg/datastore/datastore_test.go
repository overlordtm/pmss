package datastore_test

import (
	"fmt"
	"testing"

	"github.com/overlordtm/pmss/pkg/datastore"
)

func TestDatastore(t *testing.T) {
	testUris := []struct {
		testName string
		uri      string
	}{
		{
			"mysql",
			"mysql://pmss:pmss@tcp(mariadb:3306)/pmss?parseTime=true",
		},
		{
			"sqlite",
			fmt.Sprintf("sqlite3://%s/db.sqlite", t.TempDir()),
		},
	}

	for _, testcase := range testUris {
		t.Run(testcase.testName, func(t *testing.T) {
			dialector, err := datastore.ParseDBUrl(testcase.uri)
			if err != nil {
				t.Error(err)
				return
			}

			ds, err := datastore.New(datastore.WithDb(dialector))
			if err != nil {
				t.Errorf("failed to initialize datastore: %v", err)
				return
			}

			rows := []datastore.Machine{
				{
					Hostname:    "hostname1.com",
					IPv4:        "192.168.1.1",
					IPv6:        "2001:fe8::1",
					ApiKey:      "1234",
					AllowSubmit: true,
				},
				{
					Hostname:    "hostname2.com",
					IPv4:        "192.168.1.2",
					IPv6:        "2001:fe8::2",
					ApiKey:      "1234",
					AllowSubmit: true,
				},
				{
					Hostname:    "hostname3.com",
					IPv4:        "192.168.1.3",
					IPv6:        "2001:fe8::3",
					ApiKey:      "1234",
					AllowSubmit: false,
				},
			}
			if err := ds.Machines().InsertBatch(rows); err != nil {
				t.Error(err)
				return
			}
			var machine datastore.Machine
			if err := ds.Machines().FindByIPv4("192.168.1.2", &machine); err != nil {
				t.Error(err)
				return
			}
			var reportRun datastore.ReportRun
			if err := ds.ReportRuns().CreateNew(&reportRun); err != nil {
				t.Error(err)
				return
			}
			scannedFiles := []datastore.ScannedFile{
				{
					Path:      "/bin/ps",
					SHA1:      "3c738552873525fda24139e1214c95bbdaf9dcca",
					Machine:   machine,
					ReportRun: reportRun,
					Size:      137680,
					Mode:      0o755,
					Owner:     "root",
					Group:     "root",
				},
			}
			if err := ds.ScannedFiles().InsertBatch(scannedFiles); err != nil {
				t.Error(err)
				return
			}

		})
	}
}

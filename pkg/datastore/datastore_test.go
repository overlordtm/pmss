package datastore_test

import (
	"testing"

	"github.com/overlordtm/pmss/pkg/datastore"
)

func TestDatastore(t *testing.T) {
	dbPath := "mysql://pmss:pmss@tcp(mariadb:3306)/pmss"
	dialector, err := datastore.ParseDBUrl(dbPath)
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
}

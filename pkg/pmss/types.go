package pmss

import (
	"github.com/google/uuid"
	"github.com/overlordtm/pmss/pkg/datastore"
)

type ScanReport struct {
	Files     []datastore.ScannedFile
	Hostname  string
	MachineId string
	ScanRunId *uuid.UUID
}

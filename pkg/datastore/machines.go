package datastore

import "gorm.io/gorm"

// Machine represents information about a machine on the network. It also contains info whether the machine is allowed to submit files.
type Machine struct {
	ID        uint   `gorm:"primarykey"`
	Hostname  string `gorm:"type:varchar(255);uniqueIndex:;notnull"`
	MachineId string `gorm:"type:varchar(255);uniqueIndex:;notnull"`
}

type machineRepository struct {
}

func (*machineRepository) Create(machine *Machine) DbOp {
	return func(d *gorm.DB) error {
		return d.Create(machine).Error
	}
}

func (*machineRepository) CreateInBatches(batch []Machine) DbOp {
	return func(d *gorm.DB) error {
		return d.CreateInBatches(batch, 100).Error
	}
}

func (*machineRepository) FirstOrCreate(machine *Machine) DbOp {
	return func(d *gorm.DB) error {
		return d.FirstOrCreate(machine, "machine_id = ?", machine.MachineId).Error
	}
}

func (*machineRepository) FindByHostname(hostname string, outMachine *Machine) DbOp {
	return func(d *gorm.DB) error {
		return d.First(outMachine, "hostname = ?", hostname).Error
	}
}

package datastore

import (
	"gorm.io/gorm"
)

// Machine represents information about a machine on the network. It also contains info whether the machine is allowed to submit files.
type Machine struct {
	ID        uint    `gorm:"primarykey"`
	Hostname  string  `gorm:"type:varchar(255);uniqueIndex:;notnull"`
	MachineId string  `gorm:"type:varchar(255);uniqueIndex:;notnull"`
	IPv4      *string `gorm:"type:char(15)"`
	IPv6      *string `gorm:"type:char(39)"`
}

type machineRepository struct {
	db *gorm.DB
}

func (repo *machineRepository) Create(machine *Machine) error {
	return repo.db.Create(machine).Error
}

func (r *machineRepository) GetOrCreate(machine *Machine) error {
	if err := r.db.FirstOrCreate(machine, "machine_id = ?", machine.MachineId).Error; err != nil {
		return err
	}
	return nil
}

func (repo *machineRepository) FindByHostname(hostname string, outMachine *Machine) error {
	return repo.db.Find(outMachine, "hostname = ?", hostname).Error
}

func (repo *machineRepository) FindByIPv4(ipv4 string, outMachine *Machine) error {
	return repo.db.Find(outMachine, "ipv4 = ?", ipv4).Error
}

func (repo *machineRepository) FindByIPv6(ipv6 string, outMachine *Machine) error {
	return repo.db.Find(outMachine, "ipv6 = ?", ipv6).Error
}

func (repo *machineRepository) Insert(row Machine) error {
	return repo.db.Create(&row).Error
}

func (repo *machineRepository) InsertBatch(rows []Machine) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

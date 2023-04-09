package datastore

import (
	"gorm.io/gorm"
)

// Machine represents information about a machine on the network. It also contains info whether the machine is allowed to submit files.
type Machine struct {
	gorm.Model
	Hostname  string `gorm:"type:char(255);uniqueIndex:;notnull"`
	MachineId string `gorm:"type:char(255);uniqueIndex:;notnull"`
	IPv4      string `gorm:"type:char(15);uniqueIndex:;notnull"`
	IPv6      string `gorm:"type:char(39);uniqueIndex:;notnull"`
}

type machineRepository struct {
	db *gorm.DB
}

func (repo *machineRepository) Create(m *Machine) error {
	return repo.db.Create(m).Error
}

func (repo *machineRepository) FindByHostname(hostname string, dest *Machine) error {
	return repo.db.Where("hostname = ?", hostname).First(dest).Error
}

func (repo *machineRepository) FindByApiKey(api_key string, dest *Machine) error {
	return repo.db.Where("api_key = ?", api_key).First(dest).Error
}

func (repo *machineRepository) FindByIPv4(ipv4 string, dest *Machine) error {
	return repo.db.Where("ipv4 = ?", ipv4).First(dest).Error
}

func (repo *machineRepository) FindByIPv6(ipv6 string, dest *Machine) error {
	return repo.db.Where("ipv6 = ?", ipv6).First(dest).Error
}

func (repo *machineRepository) DB() *gorm.DB {
	return repo.db
}

func (repo *machineRepository) Insert(row Machine) error {
	return repo.db.Create(&row).Error
}

func (repo *machineRepository) InsertBatch(rows []Machine) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

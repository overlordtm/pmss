package datastore

import (
	"errors"

	"gorm.io/gorm"
)

// Machine represents information about a machine on the network. It also contains info whether the machine is allowed to submit files.
type Machine struct {
	Hostname    string `gorm:"primaryKey"`
	IPv4        string `gorm:"type:varchar(15),index:idx_ip4,unique"`
	IPv6        string `gorm:"type:varchar(39),index:idx_ip6,unique"`
	ApiKey      string `gorm:"type:varchar(64),index:idx_apikey,unique"`
	AllowSubmit bool   `gorm:"type:bool"`
}

type machineRepository struct {
	db *gorm.DB
}

func (repo *machineRepository) MaySubmitReports(hostname, api_key string) (bool, error) {
	var machine Machine
	failErr := errors.New("Machine not found or not allowed to submit")
	if err := repo.db.Where(&Machine{Hostname: hostname, ApiKey: api_key, AllowSubmit: true}).First(&machine).Error; err != nil {
		return false, failErr
	}
	if !machine.AllowSubmit {
		return false, failErr
	}
	return true, nil
}

func (repo *machineRepository) FindByHostname(hostname string) (*Machine, error) {
	return repo.findBy(&Machine{Hostname: hostname})
}

func (repo *machineRepository) FindByApiKey(api_key string) (*Machine, error) {
	return repo.findBy(&Machine{ApiKey: api_key})
}

func (repo *machineRepository) FindByIPv4(ipv4 string) (*Machine, error) {
	return repo.findBy(&Machine{IPv4: ipv4})
}

func (repo *machineRepository) FindByIPv6(ipv6 string) (*Machine, error) {
	return repo.findBy(&Machine{IPv6: ipv6})
}
func (repo *machineRepository) findBy(fields *Machine) (*Machine, error) {
	var row Machine
	if err := repo.db.Where(fields).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (repo *machineRepository) Insert(row Machine) error {
	return repo.db.Create(&row).Error
}

func (repo *machineRepository) InsertBatch(rows []Machine) error {
	return repo.db.CreateInBatches(&rows, 100).Error
}

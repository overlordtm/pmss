package httpserver

import "github.com/overlordtm/pmss/pkg/hashvariant"

type HashStatus string

const (
	HashStatusUnknown   HashStatus = "unknown"
	HashStatusSafe      HashStatus = "safe"
	HashStatusMalicious HashStatus = "malicious"
)

type HashResponse struct {
	Hash        string                  `json:"hash"`
	HashVariant hashvariant.HashVariant `json:"hash_variant"`
	Status      HashStatus              `json:"status"`
}

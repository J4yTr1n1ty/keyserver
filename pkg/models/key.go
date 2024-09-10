package models

import (
	"gorm.io/gorm"
)

type Key struct {
	gorm.Model
	Fingerprint string     `json:"fingerprint" gorm:"unique;not null"`
	PublicKey   string     `json:"key"`
	Identities  []Identity `json:"identities"`
}

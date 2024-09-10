package models

import (
	"time"

	"gorm.io/gorm"
)

type Key struct {
	Fingerprint string         `json:"fingerprint" gorm:"unique;not null;primaryKey"`
	PublicKey   string         `json:"key"`
	Identities  []Identity     `json:"identities" gorm:"foreignKey:KeyFingerprint;references:Fingerprint"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Key struct {
	Fingerprint string         `json:"fingerprint" gorm:"unique;not null;primaryKey"`
	PublicKey   []byte         `json:"key"`
	Identities  []Identity     `json:"identities" gorm:"foreignKey:KeyFingerprint;references:Fingerprint"`
	ValidFrom   time.Time      `json:"valid_from"`
	ValidUntil  time.Time      `json:"valid_until"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

package models

import (
	"gorm.io/gorm"
)

type Identity struct {
	gorm.Model
	Name           string `json:"name"`
	Comment        string `json:"comment"`
	Email          string `json:"email"`
	KeyFingerprint string `json:"key_fingerprint" gorm:"not null;index"` // Foreign key to Key model
}

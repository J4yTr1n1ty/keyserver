package models

import (
	"gorm.io/gorm"
)

type Identity struct {
	gorm.Model
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Email   string `json:"email"`
	KeyID   uint   `json:"key_id"`
}

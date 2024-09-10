package boot

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", Environment.GetEnv("DB_HOST"), Environment.GetEnv("DB_USER"), Environment.GetEnv("DB_PASS"), Environment.GetEnv("DB_NAME"), Environment.GetEnv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

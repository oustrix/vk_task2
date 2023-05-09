package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string, config *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), config)
}

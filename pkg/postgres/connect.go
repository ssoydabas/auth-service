package postgres

import (
	"database/sql"

	"github.com/ssoydabas/auth-service/models"
	"github.com/ssoydabas/auth-service/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPQ(config config.Config) (*gorm.DB, error) {
	sqlDB, err := sql.Open("pgx", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if ping, err := db.DB(); err != nil || ping.Ping() != nil {
		return nil, err
	}

	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

// AutoMigrate migrates the models to the database.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Account{},
		&models.AccountPassword{},
		&models.AccountToken{},
	)

}

package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mocbotau/mocbot-archive/internal/models"
)

// DB is a wrapper around the GORM database connection.
type DB struct {
	*gorm.DB
}

// NewMySQL creates a new MySQL database connection with GORM.
func NewMySQL(user, password, host, dbName string) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbInstance := &DB{db}

	err = db.AutoMigrate(&models.ListeningSession{},
		&models.TrackPlay{},
		&models.TrackPlayListener{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return dbInstance, nil
}

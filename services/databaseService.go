package services

import (
	"fmt"
	"log"
	"movie-reservation-system/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseService struct {
	DB *gorm.DB
}

func NewDatabaseService(dns string) (*DatabaseService, error) {
	fmt.Println("Initializing Database...")
	db := &DatabaseService{}
	err := db.initialize(dns)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DatabaseService) initialize(dns string) error {
	var err error

	fmt.Println("Initializing Database...")
	db.DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

func (db *DatabaseService) Migrate() error {
	err := db.DB.AutoMigrate(&models.User{}, &models.Movie{}, &models.Reservation{}, &models.Theater{}, &models.ShowTime{})
	if err != nil {
		log.Fatalf("Error auto migrate lists: %v", err)
		return err
	}

	log.Println("Migration complete")
	return nil
}

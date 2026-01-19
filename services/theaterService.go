package services

import (
	"errors"
	"fmt"
	"movie-reservation-system/models"

	"gorm.io/gorm"
)

type TheaterService struct {
	db *gorm.DB
}

func NewTheaterService(db *gorm.DB) *TheaterService {
	return &TheaterService{db: db}
}

func (ts *TheaterService) CreateTheater(theater *models.Theater) error {
	return ts.db.Create(theater).Error
}

func (ts *TheaterService) UpdateTheater(id int, theaterUpdate *models.Theater) error {
	theater, err := ts.GetTheaterByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("record with id %d not found", id)
		}
		return fmt.Errorf("ID %d retrieval error: %w", id, err)
	}

	if err := ts.db.Model(theater).Updates(theaterUpdate).Error; err != nil {
		return fmt.Errorf("failed to update theater with id %d: %w", id, err)
	}

	return nil
}

func (ts *TheaterService) DeleteTheater(id int) error {
	if err := ts.db.Delete(&models.Theater{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete theater: %w", err)
	}
	return nil
}

func (ts *TheaterService) GetTheaterByID(id int) (*models.Theater, error) {
	var theater models.Theater
	if err := ts.db.First(&theater, id).Error; err != nil {
		return nil, err
	}
	return &theater, nil
}

func (ts *TheaterService) GetAllTheaters() ([]*models.Theater, error) {
	var theaters []*models.Theater
	if err := ts.db.Find(&theaters).Error; err != nil {
		return nil, err
	}
	return theaters, nil
}

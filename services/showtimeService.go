package services

import (
	"errors"
	"fmt"
	"movie-reservation-system/models"

	"gorm.io/gorm"
)

type ShowtimeService struct {
	db *gorm.DB
}

func NewShowtimeService(db *gorm.DB) *ShowtimeService {
	return &ShowtimeService{db: db}
}

func (ss *ShowtimeService) CreateShowtime(showtime *models.ShowTime) error {
	if err := ss.db.Create(&showtime).Error; err != nil {
		return err
	}

	return nil
}

func (ss *ShowtimeService) GetShowtimeByID(id int) (*models.ShowTime, error) {
	var showtime models.ShowTime
	if err := ss.db.First(&showtime, id).Error; err != nil {
		return nil, err
	}

	return &showtime, nil
}

func (ss *ShowtimeService) UpdateShowtime(id int, updateShowtime *models.ShowTime) error {
    showtime, err := ss.GetShowtimeByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("record with id %d not found", id)
        }
        return fmt.Errorf("ID %d retrieval error: %w", id, err)
    }

    if err := ss.db.Model(showtime).Updates(updateShowtime).Error; err != nil {
        return fmt.Errorf("failed to update showtime with id %d: %w", id, err)
    }

    return nil
}


func (ss *ShowtimeService) DeleteShowtime(id uint) error {
	var showtime models.ShowTime
	if err := ss.db.First(&showtime, id).Error; err != nil {
		return err
	}

	if err := ss.db.Delete(&showtime).Error; err != nil {
		return err
	}

	return nil
}
